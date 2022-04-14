package dawg

import (
	"encoding/binary"
	"io"
)

type RankedGuide struct {
	units []RankedGuideUnit
	size  sizeType
}

func NewRankedGuide() *RankedGuide {
	return &RankedGuide{}
}

func (rg *RankedGuide) setUnits(units []GuideUnit) {
	rg.units = units
	rg.size = len(units)
}

func (rg *RankedGuide) Size() sizeType {
	return rg.size
}

func (rg *RankedGuide) TotalSize() sizeType {
	return rankedGuideUnitSize * rg.size
}

func (rg *RankedGuide) FileSize() sizeType {
	return 4 + rg.TotalSize()
}

func (rg *RankedGuide) Root() baseType {
	return 0
}

func (rg *RankedGuide) Child(index baseType) ucharType {
	return rg.units[index].Child
}

func (rg *RankedGuide) Sibling(index baseType) ucharType {
	return rg.units[index].Sibling
}

func ReadRankedGuide(r io.Reader) *RankedGuide {
	guide := NewRankedGuide()
	if !guide.Read(r) {
		return nil
	}
	return guide
}

// Reads a dictionary from an input stream.
func (rg *RankedGuide) Read(r io.Reader) bool {
	var baseSize baseType
	err := binary.Read(r, binary.LittleEndian, &baseSize)
	if err != nil {
		return false
	}

	var size sizeType = sizeType(baseSize)
	var unitsBuf = make([]GuideUnit, size)
	err = binary.Read(r, binary.LittleEndian, &unitsBuf)
	if err != nil {
		return false
	}

	rg.units = unitsBuf
	rg.size = len(unitsBuf)
	return true
}

// Writes a dictionary to an output stream.
func (rg *RankedGuide) Write(w io.Writer) bool {
	var baseSize baseType = baseType(rg.size)
	err := binary.Write(w, binary.LittleEndian, baseSize)
	if err != nil {
		return false
	}

	err = binary.Write(w, binary.LittleEndian, rg.units)
	if err != nil {
		return false
	}

	return true
}

func (rg *RankedGuide) Clear() {
	rg.setUnits(rg.units[:0])
}
