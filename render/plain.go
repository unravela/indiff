package render

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/unravela/indiff"
)

// Plain renderer is producing simple plain text output without some special structure
type Plain struct {
	rootPath          string
	showRelativePaths bool
}

// NewPlain creates new instance of Plain renderer.
// If showRelativePaths is true then printent paths will be shown relative to given rootPath.
func NewPlain(rootPath string, showRelativePaths bool) *Plain {
	return &Plain{rootPath: rootPath, showRelativePaths: showRelativePaths}
}

// Render prints given differences as simple text with one line per difference to given writer
func (p *Plain) Render(out io.Writer, diffs indiff.Diffs) {
	for _, d := range diffs {
		switch diff := d.(type) {
		case *indiff.Missing:
			fmt.Fprintf(out, "%s: missing translation of: %s\n", diff.Lang(), p.resolve(diff.Base()))
		case *indiff.ModifiedBase:
			fmt.Fprintf(out, "%s: modified only base: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
		case *indiff.ModifiedBoth:
			fmt.Fprintf(out, "%s: modified base and translation: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
		default:
			fmt.Fprintf(out, "%s: unknown difference: %s: %s\n", diff.Lang(), p.resolve(diff.Base()), p.resolve(diff.Translation()))
		}
	}
}

// resolve converts path of given file to relative path if requested and possible otherwise full path is returned
func (p *Plain) resolve(file *indiff.File) string {
	if !p.showRelativePaths {
		return file.Path
	}
	rel, err := filepath.Rel(p.rootPath, file.Path)
	if err != nil {
		return file.Path
	}
	return rel
}
