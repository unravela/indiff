package git

import (
	"bytes"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
)

// treeChanges holds changes
type treeChanges struct {
	root    string
	changes merkletrie.Changes
}

// collectChanges collects changes in repo in given revisionRange
func collectChanges(repo *git.Repository, revisionRange *Range) (*treeChanges, error) {
	// get trees for both sides of revision range
	olderTree, err := revisionRange.olderTree(repo)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get older revision tree")
	}
	newerTree, err := revisionRange.newerTree(repo)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get newer revision tree")
	}
	changes, err := merkletrie.DiffTree(olderTree, newerTree, diffTreeIsEquals)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot resolve difference between older and newer revision tree")
	}

	// resolve repo root path
	tree, err := repo.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get worktree of repository")
	}
	root := tree.Filesystem.Root()

	return &treeChanges{root: root, changes: changes}, nil
}

// forEachCreatedOrModified invokes given function f for each change with action Inserted or Modified
func (t *treeChanges) forEachCreatedOrModified(f func(path string)) {
	for _, c := range t.changes {
		action, _ := c.Action()
		if action == merkletrie.Insert || action == merkletrie.Modify {
			path := filepath.Join(t.root, c.To.String())
			f(path)
		}
	}
}

// wasCreatedOrModified check if given path points to file which was changed with action Inserted or Modified
func (t *treeChanges) wasCreatedOrModified(path string) bool {
	for _, c := range t.changes {
		action, _ := c.Action()
		if action == merkletrie.Insert || action == merkletrie.Modify {
			if path == filepath.Join(t.root, c.To.String()) {
				return true
			}
		}
	}
	return false
}

// code below was copied from go-git as it was private but needed for merkletrie.DiffTree

var emptyNoderHash = make([]byte, 24)

func diffTreeIsEquals(a, b noder.Hasher) bool {
	hashA := a.Hash()
	hashB := b.Hash()

	if bytes.Equal(hashA, emptyNoderHash) || bytes.Equal(hashB, emptyNoderHash) {
		return false
	}

	return bytes.Equal(hashA, hashB)
}
