package indiff

import (
	"reflect"
	"testing"
)

func TestBasicDiff(t *testing.T) {

	// Given bundle with "en" as base language and few files in "en" and "de"
	bundle := NewBundle("en", Files{
		NewFile("en/first.md", "en"),
		NewFile("en/second.md", "en"),
		NewFile("de/first.md", "de"),
	})

	// Given basic diff tool for "de" language
	diffTool := NewBasic([]string{"de"})

	// Given expected difference

	// When diffs are calculated
	diffs := diffTool.Diff(bundle)

	// Then diffs should contain one missing file
	var missing Diff = NewMissing(NewFile("en/second.md", "en"), "de")
	if len(diffs) != 1 {
		t.Errorf("Unexpected count of differences. Should be `%d` but was `%d`", 1, len(diffs))
	}
	if !reflect.DeepEqual(diffs[0], missing) {
		t.Errorf("Unexpected type of difference. Should be `%s` but was `%s`", missing, diffs[0])
	}
}
