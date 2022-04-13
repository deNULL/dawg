package main

import (
	"encoding/binary"
	"io"
)

type Dictionary struct {
	units []DictionaryUnit
	size  sizeType
}

func NewDictionary() *Dictionary {
	return &Dictionary{}
}

func (dict *Dictionary) setUnits(units []DictionaryUnit) {
	dict.units = units
	dict.size = len(units)
}

func (dict *Dictionary) TotalSize() sizeType {
	return unitSize * dict.size
}

func (dict *Dictionary) FileSize() sizeType {
	return 4 + dict.TotalSize()
}

// Root index.
func (dict *Dictionary) Root() baseType {
	return 0
}

// Checks if a given index is related to the end of a key.
func (dict *Dictionary) HasValue(index baseType) bool {
	return dictHasLeaf(dict.units[index])
}

// Gets a value from a given index.
func (dict *Dictionary) Value(index baseType) valueType {
	return dictValue(dict.units[index^dictOffset(dict.units[index])])
}

func ReadDictionary(r io.Reader) *Dictionary {
	dict := NewDictionary()
	if !dict.Read(r) {
		return nil
	}
	return dict
}

// Reads a dictionary from an input stream.
func (dict *Dictionary) Read(r io.Reader) bool {
	var baseSize baseType
	err := binary.Read(r, binary.LittleEndian, &baseSize)
	if err != nil {
		return false
	}

	var size sizeType = sizeType(baseSize)
	var unitsBuf = make([]DictionaryUnit, size)
	err = binary.Read(r, binary.LittleEndian, &unitsBuf)
	if err != nil {
		return false
	}

	dict.setUnits(unitsBuf)
	return true
}

// Writes a dictionary to an output stream.
func (dict *Dictionary) Write(w io.Writer) bool {
	var baseSize baseType = baseType(dict.size)
	err := binary.Write(w, binary.LittleEndian, baseSize)
	if err != nil {
		return false
	}

	err = binary.Write(w, binary.LittleEndian, dict.units)
	if err != nil {
		return false
	}

	return true
}

// Exact matching
func (dict *Dictionary) ContainsString(key string) bool {
	var index baseType = dict.Root()
	if !dict.FollowString(key, &index) {
		return false
	}
	return dict.HasValue(index)
}
func (dict *Dictionary) ContainsStringLen(key string, length sizeType) bool {
	var index baseType = dict.Root()
	if !dict.FollowStringLen(key, length, &index) {
		return false
	}
	return dict.HasValue(index)
}

// Exact matching.
func (dict *Dictionary) FindString(key string) valueType {
	var index baseType = dict.Root()
	if !dict.FollowString(key, &index) {
		return -1
	}
	if dict.HasValue(index) {
		return dict.Value(index)
	}
	return -1
}
func (dict *Dictionary) FindStringLen(key string, length sizeType) valueType {
	var index baseType = dict.Root()
	if !dict.FollowStringLen(key, length, &index) {
		return -1
	}
	if dict.HasValue(index) {
		return dict.Value(index)
	}
	return -1
}
func (dict *Dictionary) FindStringValue(key string, value *valueType) bool {
	var index baseType = dict.Root()
	if !dict.FollowString(key, &index) || !dict.HasValue(index) {
		return false
	}
	*value = dict.Value(index)
	return true
}
func (dict *Dictionary) FindStringLenValue(key string, length sizeType, value *valueType) bool {
	var index baseType = dict.Root()
	if !dict.FollowStringLen(key, length, &index) || !dict.HasValue(index) {
		return false
	}
	*value = dict.Value(index)
	return true
}

// Follows a transition.
func (dict *Dictionary) Follow(label ucharType, index *baseType) bool {
	var nextIndex baseType = *index ^ dictOffset(dict.units[*index]) ^ baseType(label)
	if dictLabel(dict.units[nextIndex]) != label {
		return false
	}
	*index = nextIndex
	return true
}

// Follows transitions.
func (dict *Dictionary) FollowString(key string, index *baseType) bool {
	for i := 0; i < len(key) && key[i] != 0; i++ {
		if !dict.Follow(key[i], index) {
			return false
		}
	}
	return true
}
func (dict *Dictionary) FollowStringCount(key string, index *baseType, count *sizeType) bool {
	for i := 0; i < len(key) && key[i] != 0; i++ {
		if !dict.Follow(key[i], index) {
			return false
		}
		*count++
	}
	return true
}

// Follows transitions.
func (dict *Dictionary) FollowStringLen(key string, length sizeType, index *baseType) bool {
	for i := 0; i < length; i++ {
		if !dict.Follow(key[i], index) {
			return false
		}
	}
	return true
}
func (dict *Dictionary) FollowStringLenCount(key string, length sizeType, index *baseType, count *sizeType) bool {
	for i := 0; i < length; i++ {
		if !dict.Follow(key[i], index) {
			return false
		}
		*count++
	}
	return true
}
