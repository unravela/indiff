package indiff

import "fmt"

// Diffs is collection of multiple differences
type Diffs = []Diff

// Diff represents same difference between base and translation file
type Diff interface {
	// Base points to base file
	Base() *File
	// Translation points to file compared with base file
	Translation() *File
	// Lang is language in which this difference occurs (usually defined by translation file)
	Lang() string
}

// DiffTool represent tool for calculating differences
type DiffTool interface {
	// Diff calculates differences
	Diff(bundle *Bundle) []Diff
}

// Missing says that there is no translation for base file in specified language
type Missing struct {
	base *File
	lang string
}

// NewMissing creates new missing file difference
func NewMissing(base *File, lang string) *Missing {
	return &Missing{base: base, lang: lang}
}

// Base points to base file for which there is no translation
func (m *Missing) Base() *File {
	return m.base
}

// Translation returns always nil because it does not exist
func (m *Missing) Translation() *File {
	return nil
}

// Lang is language in which there is no transaltion for base file
func (m *Missing) Lang() string {
	return m.lang
}

func (m *Missing) String() string {
	return fmt.Sprintf("Missing{ base: %s, lang: %s }", m.base, m.lang)
}

// Modification holds modified file changes made to this file
type Modification struct {
	file  *File
	patch string // TODO: use some structured type instead of string to be able to control output
}

// NewModification turns File to Modification (modified file) with changes in given patch
func NewModification(file *File, patch string) *Modification {
	return &Modification{file: file, patch: patch}
}

// Modified turns this file to Modification with changes in given patch
func (f *File) Modified(patch string) *Modification {
	return NewModification(f, patch)
}

// ModifiedBase says that base file was modified but it's translation was not
type ModifiedBase struct {
	base        *Modification
	translation *File
	// TODO: what exactly was added / removed
}

// NewModifiedBase creates new ModifiedBase file difference
func NewModifiedBase(base *Modification, translation *File) *ModifiedBase {
	return &ModifiedBase{base: base, translation: translation}
}

// Base points to file in base language which was modified
func (m *ModifiedBase) Base() *File {
	return m.base.file
}

// BasePatch returns text representation of changes made to file in base language
func (m *ModifiedBase) BasePatch() string {
	return m.base.patch
}

// Translation points to file which equivalent in base language was modified
func (m *ModifiedBase) Translation() *File {
	return m.translation
}

// Lang is language in which there was no modification
func (m *ModifiedBase) Lang() string {
	return m.translation.Lang
}

func (m *ModifiedBase) String() string {
	return fmt.Sprintf("ModifiedBase{ base: %s, translation: %s }", m.base.file, m.translation)
}

// ModifiedBoth says that base file and it's translation was modified
type ModifiedBoth struct {
	base        *Modification
	translation *Modification
	// TODO: what exactly was added / removed
}

// NewModifiedBoth creates new ModifiedBoth file difference
func NewModifiedBoth(base *Modification, translation *Modification) *ModifiedBoth {
	return &ModifiedBoth{base: base, translation: translation}
}

// Base points to file in base language which was modified
func (m *ModifiedBoth) Base() *File {
	return m.base.file
}

// BasePatch returns text representation of changes made to file in base language
func (m *ModifiedBoth) BasePatch() string {
	return m.base.patch
}

// TranslationPatch returns text representation of changes made to translation file
func (m *ModifiedBoth) TranslationPatch() string {
	return m.translation.patch
}

// Translation points to translation file which was modified
func (m *ModifiedBoth) Translation() *File {
	return m.translation.file
}

// Lang is language of translation file
func (m *ModifiedBoth) Lang() string {
	return m.translation.file.Lang
}

func (m *ModifiedBoth) String() string {
	return fmt.Sprintf("ModifiedBoth{ base: %s, translation: %s }", m.base.file, m.translation.file)
}
