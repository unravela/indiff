package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/unravela/indiff"
)

// Git represents diff tool based on changes in Git repository.
// Changes are bounded by specified revision range.
// It should be used only as addition to Base diff tool as it does not recognize Missing translations.
type Git struct {
	path          string
	revisionRange *Range
	repo          *git.Repository
	changes       *treeChanges
}

// ErrRepoNotFound indicates that there was no Git repository on given path
var ErrRepoNotFound = errors.New("repository not found")

// OpenGit creates Git based diff tool for repository on given path (root path or some inner path in repository) in given revisionRange
func OpenGit(path string, revisionRange *Range) (*Git, error) {
	// open repo
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		return nil, ErrRepoNotFound
	} else if err != nil {
		return nil, err
	}

	// collect changes
	changes, err := collectChanges(repo, revisionRange)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to collect changes in given range")
	}

	return &Git{path: path, revisionRange: revisionRange, repo: repo, changes: changes}, nil
}

// Diff produces differences based on changes to basefile vs changes to it's translation in specific language.
//
// It use following rules to choose Diff instance:
//
// 	Base	 	Translation		Diff
// 	-----------------------------------------
// 	insert		insert			ModifiedBoth
// 	modify		modify			ModifiedBoth
// 	delete		delete			  -
// 	insert		modify			ModifiedBoth
// 	modify		delete			  -
// 	delete		insert			  -
// 	insert		delete			  -
// 	modify		insert			ModifiedBoth
// 	delete		modify			  -
// 	insert		  -				ModifiedBase
// 	modify		  -				ModifiedBase
// 	delete		  -   			  -
//
func (g *Git) Diff(bundle *indiff.Bundle) indiff.Diffs {
	// collect only modified basepaths with their translation files
	modified := map[string][]*indiff.File{}
	g.changes.forEachCreatedOrModified(func(path string) {
		files := bundle.FilesInOtherLangs(path)
		if len(files) > 0 {
			modified[path] = files
		}
	})

	// collect translation files which were not modified or created for all modified basepaths
	diffs := []indiff.Diff{}
	for m, files := range modified {
		basefile := indiff.NewFile(m, bundle.BaseLang())
		for _, f := range files {
			// TODO: we could also produce Missing difference here
			wasTranslationModified := g.changes.wasCreatedOrModified(f.Path)
			if wasTranslationModified {
				diffs = append(diffs, indiff.NewModifiedBoth(basefile, f))
			} else {
				diffs = append(diffs, indiff.NewModifiedBase(basefile, f))
			}
		}
	}

	return diffs
}
