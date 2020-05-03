package indiff

// Bundle holds files in base language and corresponsing translation files in different languages
type Bundle struct {
	baselang        string
	filesByLang     map[string]Files
	filesByBasepath map[string]map[string]*File
}

// NewBundle creates bundle with for specified baselang and files collection
func NewBundle(baselang string, files Files) *Bundle {
	filesByLang := map[string]Files{}
	for _, f := range files {
		filesByLang[f.Lang] = append(filesByLang[f.Lang], f)
	}

	filesByBasePath := make(map[string]map[string]*File, len(filesByLang[baselang]))
	for _, bf := range filesByLang[baselang] {

		filesByBasePath[bf.Path] = map[string]*File{}

		for lang, files := range filesByLang {
			// skip base lang to not check "against self"
			if lang == baselang {
				continue
			}

			for _, f := range files {
				if bf.IsEqualInOtherLang(f) {
					filesByBasePath[bf.Path][f.Lang] = f
					break
				}
			}
		}
	}

	return &Bundle{
		baselang:        baselang,
		filesByLang:     filesByLang,
		filesByBasepath: filesByBasePath,
	}
}

// BaseLang returns base language of bundle
func (b *Bundle) BaseLang() string {
	return b.baselang
}

// Langs lists all languages in bundle including base language
func (b *Bundle) Langs() []string {
	langs := make([]string, len(b.filesByLang))
	for l := range b.filesByLang {
		langs = append(langs, l)
	}
	return langs
}

// BaseFiles returns all files in base language
func (b *Bundle) BaseFiles() Files {
	files := Files{}
	for p := range b.filesByBasepath {
		files = append(files, NewFile(p, b.baselang))
	}
	return files
}

// BasePaths returns paths to all files in base language
func (b *Bundle) BasePaths() []string {
	paths := []string{}
	for p := range b.filesByBasepath {
		paths = append(paths, p)
	}
	return paths
}

// FilesForLang returns all files in specified language
func (b *Bundle) FilesForLang(lang string) Files {
	return b.filesByLang[lang]
}

// FileInLang returns single file with translation of file specified by basepath in given language
func (b *Bundle) FileInLang(basepath string, lang string) *File {
	return b.filesByBasepath[basepath][lang]
}

// FilesInOtherLangs returns all files with transaltion of file specified by basepath
func (b *Bundle) FilesInOtherLangs(basepath string) Files {
	section := b.filesByBasepath[basepath]
	if section == nil {
		return nil
	}
	files := Files{}
	for _, f := range section {
		files = append(files, f)
	}
	return files
}
