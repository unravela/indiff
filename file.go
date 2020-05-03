package indiff

import (
	"fmt"
	"strings"
)

// File is just some path to translation file in some language
type File struct {
	// Path to translation file
	Path string
	// Lang as language of translation in this file
	Lang string
}

// Files represent collection of File
type Files = []*File

// NewFile creates new file for specified path and language
func NewFile(path string, lang string) *File {
	return &File{Path: path, Lang: lang}
}

// IsEqualInOtherLang checks if given other file is translation of this file in other language
func (f *File) IsEqualInOtherLang(other *File) bool {
	// we expect language codes have same length, so paths should have same length too
	if len(f.Path) != len(other.Path) {
		return false
	}

	// create diff string with spaces for common characters and with characters from base path for conflict characters
	ar := []rune(f.Path)
	br := []rune(other.Path)
	cr := []rune(strings.Repeat(" ", len(f.Path)))
	for i := 0; i < len(ar); i++ {
		if ar[i] != br[i] {
			cr[i] = ar[i]
		}
	}

	// remove spaces around difference to potentionally get only language code as difference
	diff := strings.Trim(string(cr), " ")

	// there are some additional characters, not only language code
	if len(diff) > len(f.Lang) {
		return false
	}

	// contains solves cases when language codes have some common characters (e.g. `sk` as base and `sl` will have diff only `k`)
	return strings.Contains(f.Lang, diff)
}

func (f *File) String() string {
	return fmt.Sprintf("File{ path: %s, lang: %s }", f.Path, f.Lang)
}
