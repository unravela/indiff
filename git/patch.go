package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/utils/diff"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	fdiff "github.com/go-git/go-git/v5/plumbing/format/diff"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

// revisionPatch holds multiple file patches with changes between revisions
// It implements interface Patch defined in github.com/go-git/go-git/v5/plumbing/format/diff
// It was needed to re-implement this interface because of go-git implementation is private and it does not support changes from worktree.
type revisionPatch struct {
	message     string
	filePatches []fdiff.FilePatch
}

// createPatchFromSingleChange turns one revisionChange into revisionPatch with one filePatch
func createPatchFromSingleChange(c *revisionChange) *revisionPatch {
	diffs := diff.Do(c.fromContent(), c.toContent())

	var chunks []fdiff.Chunk
	for _, d := range diffs {

		var op fdiff.Operation
		switch d.Type {
		case dmp.DiffEqual:
			op = fdiff.Equal
		case dmp.DiffDelete:
			op = fdiff.Delete
		case dmp.DiffInsert:
			op = fdiff.Add
		}

		chunks = append(chunks, &textChunk{content: d.Text, operation: op})
	}

	fp := &filePatch{
		chunks: chunks,
	}
	if c.fromPath() != "" {
		fp.from = &regularFile{c.fromPath()}
	}
	if c.toPath() != "" {
		fp.to = &regularFile{c.toPath()}
	}

	return &revisionPatch{
		filePatches: []fdiff.FilePatch{fp},
	}
}

// String returns text representation of patch.
// It does not contain header information only chunks.
func (p *revisionPatch) String() string {
	builder := &strings.Builder{}
	encoder := fdiff.NewUnifiedEncoder(builder, fdiff.DefaultContextLines)
	encoder.Encode(p)
	return removeHeader(strings.TrimSpace(builder.String()))
}

func removeHeader(patch string) string {
	split := strings.SplitN(patch, "@@", 2)
	if len(split) < 2 {
		return ""
	}
	return fmt.Sprintf("@@%s", split[1])
}

// FilePatches returns all file patches in this patch
func (p *revisionPatch) FilePatches() []fdiff.FilePatch {
	return p.filePatches
}

// Message returns additional message attached to this patch
func (p *revisionPatch) Message() string {
	return p.message
}

// filePatch holds chunks for one file
// It implements interface FilePatch defined in github.com/go-git/go-git/v5/plumbing/format/diff
type filePatch struct {
	isBinary bool
	from, to fdiff.File
	chunks   []fdiff.Chunk
}

// IsBinary returns whether file is binary
func (f *filePatch) IsBinary() bool {
	return f.isBinary
}

// Files returns file information for before and after patch versions
func (f *filePatch) Files() (from, to fdiff.File) {
	return f.from, f.to
}

// Chunks returns all chunks from this file
func (f *filePatch) Chunks() []fdiff.Chunk {
	return f.chunks
}

// textChunk represent one change in textFile
// It implements interface textChunk defined in github.com/go-git/go-git/v5/plumbing/format/diff
type textChunk struct {
	content   string
	operation fdiff.Operation
}

// Content returns text changes
func (c *textChunk) Content() string {
	return c.content
}

// Type returns type of operation in this chunk
func (c *textChunk) Type() fdiff.Operation {
	return c.operation
}

// regularFile represents non-executable file without hash
// It implements interface textChunk defined in github.com/go-git/go-git/v5/plumbing/format/diff
type regularFile struct {
	path string
}

// Hash always returns zero hash
func (f *regularFile) Hash() plumbing.Hash {
	return plumbing.ZeroHash
}

// Mode always returns regular
func (f *regularFile) Mode() filemode.FileMode {
	return filemode.Regular
}

// Path returns path to file relative to git repository root
func (f *regularFile) Path() string {
	return f.path
}
