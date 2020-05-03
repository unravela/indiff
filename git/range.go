package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie/filesystem"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
)

// Range represents revision range used to look up for the changes.
// It's defined by Older and Newer revision names, where name can be commit, tag or reference.
// Tilde and carret are supported in names too.
// If Older is not defined, HEAD is used as default.
// If Newer is not defined, working tree with staged and untracked files is used.
type Range struct {
	// Older contains name of revision from which range starts
	Older string
	// Newer contains name of revision where range ends
	Newer string
}

// NewRange creates new revision range.
// You can use empty string for older and newer arguments to use default values (see Range description).
func NewRange(older string, newer string) *Range {
	return &Range{Older: older, Newer: newer}
}

// Uncommited is default range and identifies only not yet commited (staged + untracked)
var Uncommited = &Range{}

// olderTree returns tree for older revision in given repo
func (r *Range) olderTree(repo *git.Repository) (noder.Noder, error) {
	rev := "HEAD"
	if r.Older != "" {
		rev = r.Older
	}
	return revisionTree(repo, rev)
}

// olderTree returns tree for newer revision in given repo
func (r *Range) newerTree(repo *git.Repository) (noder.Noder, error) {
	if r.Newer == "" {
		return workingTree(repo)
	}
	return revisionTree(repo, r.Newer)
}

// revisionTree returns tree for given repo and rev as revision name
func revisionTree(repo *git.Repository, rev string) (noder.Noder, error) {
	hash, err := repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return nil, err
	}
	commit, err := repo.CommitObject(*hash)
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	return object.NewTreeRootNode(tree), nil
}

// workingTree returns working tree for given repo
func workingTree(repo *git.Repository) (noder.Noder, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	var submodules map[string]plumbing.Hash // = TODO: add support for submodules ?
	// currently hidden as private function on Worktree: getSubmodulesStatus()
	// see: https://github.com/go-git/go-git/blob/e04168bb11a960018b6bbabd6972fd33163b6f28/worktree_status.go#L179
	return filesystem.NewRootNode(tree.Filesystem, submodules), nil
}
