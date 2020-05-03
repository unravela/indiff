package filesystem

import (
	"fmt"
	"os"
	"strings"

	"github.com/gobwas/glob"
)

// Pattern represent GLOB like pattern for matching paths with translation files
type Pattern string

// PredefinedPatterns contains named patterns.
// Key is the name of pattern, value array contains pattern as first element and description as second one.
var PredefinedPatterns = map[string][]string{
	"SUB": {"%l/**.%e", "each language in separate subdirectory"},
	"EXT": {"**.%l.%e", "language code as part of file extension"},
}

// ParsePattern validates given rawPattern and apply given extensions to create new Pattern.
// Given rawPattern must contain `%l` placeholder, which will be replaced later by specific language code.
// Given rawPattern can contain `%e` placeholder, which will be replaced by given extensions.
// Instead of rawPattern you can also provide one of the keys from PredefinedPatterns.
// If you provide empty extensions, `*` will be used to match any extension.
//
// E.g. ParsePattern("%l/**.%e", {"md","rst"}) will produce pattern "%l/**.{md,rst}"
func ParsePattern(rawPattern string, extensions []string) (Pattern, error) {
	// parse pattern
	pattern := PredefinedPatterns[rawPattern][0]
	if pattern == "" {
		pattern = rawPattern
	}
	if !strings.Contains(pattern, "%l") {
		return "", fmt.Errorf("Pattern must contain placeholder for language code '%%l'")
	}

	// parse extensions
	extpattern := "*"
	if len(extensions) > 0 {
		extpattern = "{" + strings.Join(extensions, ",") + "}"
	}

	// apply extensions to pattern
	pattern = strings.Replace(pattern, "%e", extpattern, 1)
	
	return Pattern(pattern), nil
}

// MustParsePattern can be used when you are sure that your pattern is valid (e.g. one of predefined patterns).
// It does same parsing as ParsePattern but panics when there will be an error.
func MustParsePattern(rawPattern string, extensions []string) Pattern {
	pattern, err := ParsePattern(rawPattern, extensions)
	if err != nil {
		panic(err)
	}
	return pattern
}

// Compile turns pattern into matcher for given lang
func (p Pattern) Compile(lang string) glob.Glob {
	rawglob := strings.Replace(string(p), "%l", lang, 1)
	return glob.MustCompile(rawglob, os.PathSeparator)
}