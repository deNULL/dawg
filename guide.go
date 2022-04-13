package main

import (
	"encoding/binary"
	"io"
)

type Guide struct {
	units []GuideUnit
	size  sizeType
}

func NewGuide() *Guide {
	return &Guide{}
}

func (guide *Guide) setUnits(units []GuideUnit) {
	guide.units = units
	guide.size = len(units)
}

func (guide *Guide) TotalSize() sizeType {
	return guideUnitSize * guide.size
}

func (guide *Guide) FileSize() sizeType {
	return 4 + guide.TotalSize()
}

// The root index
func (guide *Guide) Root() baseType {
	return 0
}

func (guide *Guide) Child(index baseType) ucharType {
	return guide.units[index].Child
}

func (guide *Guide) Sibling(index baseType) ucharType {
	return guide.units[index].Sibling
}

func ReadGuide(r io.Reader) *Guide {
	guide := NewGuide()
	if !guide.Read(r) {
		return nil
	}
	return guide
}

// Reads a dictionary from an input stream.
func (guide *Guide) Read(r io.Reader) bool {
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

	guide.units = unitsBuf
	guide.size = len(unitsBuf)
	return true
}

// Writes a dictionary to an output stream.
func (guide *Guide) Write(w io.Writer) bool {
	var baseSize baseType = baseType(guide.size)
	err := binary.Write(w, binary.LittleEndian, baseSize)
	if err != nil {
		return false
	}

	err = binary.Write(w, binary.LittleEndian, guide.units)
	if err != nil {
		return false
	}

	return true
}

func (guide *Guide) Clear() {
	guide.setUnits(guide.units[:0])
}
