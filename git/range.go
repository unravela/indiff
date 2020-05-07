package git

import (
	"fmt"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie/filesystem"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
	"github.com/pkg/errors"
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

// Uncommited is default range and identifies only not yet commited (staged + untracked)
var Uncommited = &Range{}

// revisionRange is already resolved Range
type revisionRange struct {
	ref   *Range
	older *revisionTree
	newer *revisionTree
}

// revisionTree enables access to files at specific revision
type revisionTree struct {
	root      noder.Noder
	contentOf func(path string) (string, error)
}

// ErrFileNotFound is returned when requesting file from revision tree which is not part of tree
var ErrFileNotFound = errors.New("File not found")

// resolve range to revisionRange validates and opens "both sides" of the range
func (r *Range) resolve(repo *git.Repository) (*revisionRange, error) {
	older, err := resolveOlderTree(repo, r)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid older revision")
	}
	newer, err := resolveNewerTree(repo, r)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid newer revision")
	}
	return &revisionRange{ref: r, older: older, newer: newer}, nil
}

// resolveOlderTree returns tree for older revision in given repo
func resolveOlderTree(repo *git.Repository, r *Range) (*revisionTree, error) {
	rev := "HEAD"
	if r.Older != "" {
		rev = r.Older
	}
	return commitTree(repo, rev)
}

// resolveNewerTree returns tree for newer revision in given repo
func resolveNewerTree(repo *git.Repository, r *Range) (*revisionTree, error) {
	if r.Newer == "" {
		return workingTree(repo)
	}
	return commitTree(repo, r.Newer)
}

// commitTree returns revisionTree for given repo and revision string resolvable to commit hash
func commitTree(repo *git.Repository, rev string) (*revisionTree, error) {
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
	return &revisionTree{
		root: object.NewTreeRootNode(tree),
		contentOf: func(path string) (string, error) {
			f, err := commit.File(path)
			if err != nil {
				return "", fmt.Errorf("File not found in revision %s: %s", path, rev)
			}
			c, err := f.Contents()
			if err != nil {
				return "", errors.Wrapf(err, "Unable to read content of file from revision %s: %s", path, rev)
			}
			return c, nil
		},
	}, nil
}

// workingTree returns revisionTree based on working tree for given repo
func workingTree(repo *git.Repository) (*revisionTree, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	var submodules map[string]plumbing.Hash // = TODO: add support for submodules ?
	// currently hidden as private function on Worktree: getSubmodulesStatus()
	// see: https://github.com/go-git/go-git/blob/e04168bb11a960018b6bbabd6972fd33163b6f28/worktree_status.go#L179
	return &revisionTree{
		root: filesystem.NewRootNode(tree.Filesystem, submodules),
		contentOf: func(path string) (string, error) {
			f, err := tree.Filesystem.Open(path)
			if err != nil {
				return "", fmt.Errorf("File not found in working tree: %s", path)
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return "", errors.Wrapf(err, "Unable to read content of file from working tree: %s", path)
			}
			return string(b), nil
		},
	}, nil
}
