package git

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
)

// revChange holds one change between to revisions
type revisionChange struct {
	underlying    merkletrie.Change
	revisionRange *revisionRange
}

// revisionChanges represent collection of changes between two revisions
type revisionChanges []*revisionChange

// collectChanges collects changes in repo in given revisionRange
func collectChanges(repo *git.Repository, revisionRange *revisionRange) (revisionChanges, error) {
	older := revisionRange.older.root
	newer := revisionRange.newer.root
	originalChanges, err := merkletrie.DiffTree(older, newer, diffTreeIsEquals)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot resolve difference between older and newer revision tree")
	}

	// enrich changes with revision references
	changes := make(revisionChanges, len(originalChanges))
	for i, c := range originalChanges {
		changes[i] = &revisionChange{
			underlying:    c,
			revisionRange: revisionRange,
		}
	}

	return changes, nil
}

// forEachCreatedOrModified invokes given function f for each change with action Inserted or Modified
func (changes revisionChanges) forEachCreatedOrModified(f func(change *revisionChange)) {
	for _, c := range changes {
		action, _ := c.underlying.Action()
		if action == merkletrie.Insert || action == merkletrie.Modify {
			f(c)
		}
	}
}

// fromPath returns path to file before change
func (c *revisionChange) fromPath() string {
	return c.underlying.From.String()
}

// fromContent reads content of the file before change
func (c *revisionChange) fromContent() string {
	if c.fromPath() == "" {
		return ""
	}
	content, err := c.revisionRange.older.contentOf(c.fromPath())
	if err != nil {
		panic(err)
	}
	return content
}

// toPath returns path to file after change
func (c *revisionChange) toPath() string {
	return c.underlying.To.String()
}

// toContent returns content of the file after change
func (c *revisionChange) toContent() string {
	if c.toPath() == "" {
		return ""
	}
	content, err := c.revisionRange.newer.contentOf(c.toPath())
	if err != nil {
		panic(err)
	}
	return content
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
