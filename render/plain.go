package render

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/unravela/indiff"
)

// Plain renderer is producing simple plain text output without some special structure
type Plain struct {
	RootPath          string
	ShowRelativePaths bool
	ShowDiff          bool
}

// Render prints given differences as simple text with one line per difference to given writer
func (p *Plain) Render(out io.Writer, diffs indiff.Diffs) {
	for _, d := range diffs {
		switch diff := d.(type) {
		case *indiff.Missing:
			fmt.Fprintf(out, "%s: missing translation of: %s\n", diff.Lang(), p.resolve(diff.Base()))
		case *indiff.ModifiedBase:
			fmt.Fprintf(out, "%s: modified only base: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
			p.renderDiff(out, diff.Base(), diff.BasePatch())
		case *indiff.ModifiedBoth:
			fmt.Fprintf(out, "%s: modified base and translation: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
			p.renderDiff(out, diff.Base(), diff.BasePatch())
			p.renderDiff(out, diff.Translation(), diff.TranslationPatch())
		default:
			fmt.Fprintf(out, "%s: unknown difference: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
		}
	}
}

// renderDiff prints changes made in file
func (p *Plain) renderDiff(out io.Writer, f *indiff.File, patch string) {
	if p.ShowDiff {
		patchText := patch
		if patchText == "" {
			patchText = "<no content>"
		}
		fmt.Fprintf(out, "diff %s\n%s\n", p.resolve(f), patchText)
	}
}

// resolve converts path of given file to relative path if requested and possible otherwise full path is returned
func (p *Plain) resolve(file *indiff.File) string {
	if !p.ShowRelativePaths {
		return file.Path
	}
	rel, err := filepath.Rel(p.RootPath, file.Path)
	if err != nil {
		return file.Path
	}
	return rel
}
