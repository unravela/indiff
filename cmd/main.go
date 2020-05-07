package main

import (
	"os"

	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/unravela/indiff"
	"github.com/unravela/indiff/filesystem"
	"github.com/unravela/indiff/git"
	"github.com/unravela/indiff/render"
)

func main() {

	app := &cli.App{
		Name:      "indiff",
		Usage:     "looks for missing transaltions",
		ArgsUsage: "languages",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "baselang",
				Usage:       "Base language `CODE` against which diffs in other languages are tested",
				Aliases:     []string{"b"},
				DefaultText: "first provided language",
			},
			&cli.StringFlag{
				Name: "glob",
				Usage: trimMargin(
					"Glob `PATTERN` " + `for language files identification.
					|	You can use following placeholders in pattern: 
					|		- %l: language code (required)
					|		- %e: one or more supported file extensions
					|	You can also use one of predefined patterns: ` + listPredefinedPatterns()),
				Aliases: []string{"g"},
				Value:   "SUB",
			},
			&cli.StringFlag{
				Name:        "directory",
				Usage:       "Working directory `PATH`",
				Aliases:     []string{"d"},
				Value:       ".",
				DefaultText: "current dir",
			},
			&cli.StringSliceFlag{
				Name:    "extensions",
				Usage:   "File extensions substituted for %e in given glob (default: any)",
				Aliases: []string{"e"},
			},
			&cli.BoolFlag{
				Name:  "no-git",
				Usage: "Do not use Git",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "from-revision",
				Usage:   "Revision in Git repository from which changes will be calculated (default: HEAD)",
				Aliases: []string{"f"},
			},
			&cli.StringFlag{
				Name:    "to-revision",
				Usage:   "Revision in Git repository to which changes will be calculated (default: changes in worktree)",
				Aliases: []string{"t"},
			},
			&cli.BoolFlag{
				Name:    "absolute-paths",
				Usage:   "Print absolute paths",
				Value:   false,
				Aliases: []string{"a"},
			},
			&cli.BoolFlag{
				Name:    "show-diff",
				Usage:   "Print diff for each modified file",
				Value:   false,
				Aliases: []string{"i"},
			},
		},
		Writer:          os.Stderr,
		HideHelpCommand: true,
		Action:          run,
	}

	// run application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	// parse langs
	rawlangs := c.Args().First()
	// TODO: add autodiscovery of languages if one of predefined globs used
	if rawlangs == "" {
		cli.ShowAppHelp(c)
		return fmt.Errorf("Missing required argument: languages")
	}
	langs := strings.Split(rawlangs, ",")
	if len(langs) < 2 {
		cli.ShowAppHelp(c)
		return fmt.Errorf("Invalid argument: languages: provide minimally two language codes separated by comma")
	}

	// parse baselang
	baselang := c.String("baselang")
	if baselang == "" {
		baselang = langs[0]
	} else if !contains(langs, baselang) {
		cli.ShowAppHelp(c)
		return fmt.Errorf("Invalid argument: baselang: language '%s' not found", baselang)
	}

	// parse working direcotry
	root, err := filepath.Abs(c.String("directory"))
	if err != nil {
		cli.ShowAppHelp(c)
		return errors.Wrap(err, "Invalid argument: directory")
	}

	// parse pattern
	pattern, err := filesystem.ParsePattern(c.String("glob"), c.StringSlice("extensions"))
	if err != nil {
		cli.ShowAppHelp(c)
		return errors.Wrap(err, "Invalid argument: glob")
	}

	// parse revision range
	revisionRange := &git.Range{
		Older: c.String("from-revision"),
		Newer: c.String("to-revision"),
	}
	isGitAllowed := !c.Bool("no-git")

	// collect bundle
	files := filesystem.NewFs(root, pattern).CollectFiles(langs)
	bundle := indiff.NewBundle(baselang, files)

	// calculate basic diffs
	diffs := indiff.NewBasic(langs).Diff(bundle)

	// calculate git based diffs
	if isGitAllowed {
		g, err := git.OpenGit(root, revisionRange)
		if err == git.ErrRepoNotFound {
			fmt.Fprintf(os.Stderr, "WARN: Git repository was not found. Check your path or use --no-git to hide this warning.\n")
			return nil
		} else if err != nil {
			return errors.Wrap(err, "Error during opening Git repository")
		}
		diffs = append(diffs, g.Diff(bundle)...)
	}

	// render
	r := &render.Plain{
		RootPath:          root,
		ShowRelativePaths: !c.Bool("absolute-paths"),
		ShowDiff:          c.Bool("show-diff"),
	}
	r.Render(os.Stdout, diffs)

	return nil
}

// helpers

func contains(xs []string, x string) bool {
	for _, s := range xs {
		if s == x {
			return true
		}
	}
	return false
}

func listPredefinedPatterns() string {
	var sb strings.Builder
	for id, p := range filesystem.PredefinedPatterns {
		fmt.Fprintf(&sb, "\n\t\t- %s: %s (%s)", id, p[1], p[0])
	}
	return sb.String()
}

var trimMarginRegexp = regexp.MustCompile("(?m)^\\s*\\|")

func trimMargin(s string) string {
	return trimMarginRegexp.ReplaceAllString(s, "")
}
