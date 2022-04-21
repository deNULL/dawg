package dawg

import (
	"encoding/binary"
	"io"
)

type IndexUnit baseType

const indexUnitSize = 4

type Index struct {
	units []IndexUnit
}

func NewIndex() *Index {
	return &Index{}
}

func (index *Index) Size() sizeType {
	return len(index.units)
}

func (index *Index) TotalSize() sizeType {
	return indexUnitSize * len(index.units)
}

func (index *Index) FileSize() sizeType {
	return 4 + index.TotalSize()
}

// The root index
func (index *Index) Root() baseType {
	return 0
}

func (index *Index) ChildCount(i baseType) baseType {
	return baseType(index.units[i])
}

func (index *Index) TotalCount() baseType {
	return baseType(index.units[0])
}

func ReadIndex(r io.Reader) *Index {
	index := NewIndex()
	if !index.Read(r) {
		return nil
	}
	return index
}

// Reads an index from an input stream.
func (index *Index) Read(r io.Reader) bool {
	var baseSize baseType
	err := binary.Read(r, binary.LittleEndian, &baseSize)
	if err != nil {
		return false
	}

	var size sizeType = sizeType(baseSize)
	index.units = make([]IndexUnit, size)
	err = binary.Read(r, binary.LittleEndian, &index.units)
	if err != nil {
		return false
	}
	return true
}

// Writes an index to an output stream.
func (index *Index) Write(w io.Writer) bool {
	var baseSize baseType = baseType(len(index.units))
	err := binary.Write(w, binary.LittleEndian, baseSize)
	if err != nil {
		return false
	}

	err = binary.Write(w, binary.LittleEndian, index.units)
	if err != nil {
		return false
	}

	return true
}
