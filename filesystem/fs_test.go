package filesystem

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/unravela/indiff"
)

func TestCollectFiles(t *testing.T) {

	t.Run("SUB pattern", func(t *testing.T) {
		// Given folder with translation files
		root := rootFolder("sub")

		// Given files collector with predefined pattern for "subdirectory layout" and markdown files
		fs := NewFs(root, MustParsePattern("SUB", []string{"md"}))

		// When files are collected for languages "en" and "de"
		files := fs.CollectFiles([]string{"en", "de"})

		// Then expected files should be found in collection
		expected := indiff.Files{
			indiff.NewFile(filepath.Join(root, "en", "first.md"), "en"),
			indiff.NewFile(filepath.Join(root, "en", "second.md"), "en"),
			indiff.NewFile(filepath.Join(root, "en", "section", "one.md"), "en"),
			indiff.NewFile(filepath.Join(root, "en", "section", "two.md"), "en"),
			indiff.NewFile(filepath.Join(root, "de", "first.md"), "de"),
			indiff.NewFile(filepath.Join(root, "de", "section", "one.md"), "de"),
		}
		assertCollected(t, expected, files)
	})

	t.Run("EXT pattern", func(t *testing.T) {
		// Given folder with translation files
		root := rootFolder("ext")

		// Given files collector with predefined pattern for "subdirectory layout" and markdown files
		fs := NewFs(root, MustParsePattern("EXT", []string{"md"}))

		// When files are collected for languages "en" and "de"
		files := fs.CollectFiles([]string{"en", "de"})

		// Then expected files should be found in collection
		expected := indiff.Files{
			indiff.NewFile(filepath.Join(root, "first.en.md"), "en"),
			indiff.NewFile(filepath.Join(root, "second.en.md"), "en"),
			indiff.NewFile(filepath.Join(root, "section", "one.en.md"), "en"),
			indiff.NewFile(filepath.Join(root, "section", "two.en.md"), "en"),
			indiff.NewFile(filepath.Join(root, "first.de.md"), "de"),
			indiff.NewFile(filepath.Join(root, "section", "one.de.md"), "de"),
		}
		assertCollected(t, expected, files)
	})

}

// helpers

func rootFolder(dir string) string {
	root := filepath.Join("..", "testdata", "layout", dir)
	root, _ = filepath.Abs(root)
	return root
}

func assertCollected(t *testing.T, expected, collected indiff.Files) {
	if !reflect.DeepEqual(expected, collected) {
		t.Errorf("Unexpected files collected.\n\nExpected: \n%s \n\nCollected: \n%s", printFiles(expected), printFiles(collected))
	}
}

func printFiles(files indiff.Files) string {
	b := &strings.Builder{}
	for _, f := range files {
		fmt.Fprintf(b, "%s\n", f)
	}
	return b.String()
}
