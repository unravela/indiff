package indiff

// Basic represents diff tool which is able to calculate differences only from file paths and languages.
// It does not use any other external tool and does not need to read the contents of the files.
type Basic struct {
	langs []string
}

// NewBasic creates new base diff tool for specified languages
func NewBasic(langs []string) *Basic {
	return &Basic{langs: langs}
}

// Diff calculates the differences in given bundle
func (b *Basic) Diff(bundle *Bundle) []Diff {
	diffs := []Diff{}
	for _, lang := range b.langs {
		if lang == bundle.BaseLang() {
			// skip base lang to not check "against self"
			continue
		}
		for _, basefile := range bundle.BaseFiles() {
			if bundle.FileInLang(basefile.Path, lang) == nil {
				diffs = append(diffs, NewMissing(basefile, lang))
			}
		}
	}
	return diffs
}
