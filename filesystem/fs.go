package filesystem

import (
	"os"
	"path/filepath"

	"github.com/unravela/indiff"
)

// Fs collects translation files from filesystem
type Fs struct {
	root    string
	pattern Pattern
}

// NewFs creates new instance of Fs under specified root directory.
// It accpets also GLOB like pattern which is used to distinguish which path belong to which language.
func NewFs(root string, pattern Pattern) *Fs {
	return &Fs{
		root:    root,
		pattern: pattern,
	}
}

// CollectFiles collects Files from file system for specified langs
func (fs *Fs) CollectFiles(langs []string) indiff.Files {
	files := []*indiff.File{}
	for _, lang := range langs {
		glob := fs.pattern.Compile(lang)

		filepath.Walk(fs.root, func(path string, info os.FileInfo, err error) error {
			rel, _ := filepath.Rel(fs.root, path)
			abs, _ := filepath.Abs(path)
			if glob.Match(rel) {
				files = append(files, indiff.NewFile(abs, lang))
			}
			return nil
		})
	}
	return files
}
