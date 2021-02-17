package editor

import (
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// RenameBlock reads HCL from io.Reader, and renames matched blocks, and writes
// the updated HCL to io.Writer.
// Note that a filename is used only for an error message.
// If an error occurs, Nothing is written to the output stream.
func RenameBlock(r io.Reader, w io.Writer, filename string, from string, to string) error {
	filter := &blockRename{from: from, to: to}
	return FilterHCL(r, w, filename, filter)
}

// blockRename is a filter implementation for renaming block.
type blockRename struct {
	from string
	to   string
}

// Filter reads HCL and renames matched blocks at a given address.
// The blocks which do not match the from address are output as is.
// Rename means setting the block type and labels corresponding to the new
// address.  changing the block type does not make sense on an application
// context, but filters can chain to others and the later filter may edit its
// attributes. So we allow this filter to any block type and labels.
func (f *blockRename) Filter(inFile *hclwrite.File) (*hclwrite.File, error) {
	fromTypeName, fromLabels, err := parseAddress(f.from)
	if err != nil {
		return nil, err
	}

	toTypeName, toLabels, err := parseAddress(f.to)
	if err != nil {
		return nil, err
	}

	matched := findBlocks(inFile.Body(), fromTypeName, fromLabels)

	for _, b := range matched {
		b.SetType(toTypeName)
		b.SetLabels(toLabels)
	}

	return inFile, nil
}
