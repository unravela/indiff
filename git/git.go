package git

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/unravela/indiff"
)

// Git represents diff tool based on changes in Git repository.
// Changes are bounded by specified revision range.
// It should be used only as addition to Base diff tool as it does not recognize Missing translations.
type Git struct {
	path          string
	revisionRange *revisionRange
	repo          *git.Repository
	changes       revisionChanges
}

// ErrRepoNotFound indicates that there was no Git repository on given path
var ErrRepoNotFound = errors.New("repository not found")

// OpenGit creates Git based diff tool for repository on given path (root path or some inner path in repository) in given revisionRange
func OpenGit(path string, rangeRef *Range) (*Git, error) {
	// open repo
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		return nil, ErrRepoNotFound
	} else if err != nil {
		return nil, err
	}

	// resolve repo root path
	tree, err := repo.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get worktree of repository")
	}
	rootPath := tree.Filesystem.Root()

	// resolve Range
	revisionRange, err := rangeRef.resolve(repo)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to resolve revision range")
	}

	// collect changes
	changes, err := collectChanges(repo, revisionRange)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to collect changes in given range")
	}

	return &Git{path: rootPath, revisionRange: revisionRange, repo: repo, changes: changes}, nil
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
	// collect only modified changes
	modified := map[string]*revisionChange{}
	g.changes.forEachCreatedOrModified(func(change *revisionChange) {
		path := filepath.Join(g.path, change.toPath())
		modified[path] = change
	})

	// create diffs from modified changes on files in base language
	diffs := []indiff.Diff{}
	for path, baseChange := range modified {
		files := bundle.FilesInOtherLangs(path)
		if len(files) > 0 {
			base := modify(indiff.NewFile(path, bundle.BaseLang()), baseChange)
			for _, f := range files {
				fileChange := modified[f.Path]
				if fileChange != nil {
					diffs = append(diffs, indiff.NewModifiedBoth(base, modify(f, fileChange)))
				} else {
					diffs = append(diffs, indiff.NewModifiedBase(base, f))
				}
			}
		}
	}

	return diffs
}

func modify(file *indiff.File, c *revisionChange) *indiff.Modification {
	// TODO: patch is calculated even if it is not shown at the end, find some "lazy loading" solution
	patch := createPatchFromSingleChange(c)
	return file.Modified(patch.String())
}
