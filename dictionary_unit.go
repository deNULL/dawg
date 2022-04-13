package main

const offsetMax = baseType(1) << 21
const isLeafBit = baseType(1) << 31
const hasLeafBit = baseType(1) << 8
const extensionBit = baseType(1) << 9

type DictionaryUnit = baseType

const unitSize = 4 // Replacement for sizeof(DictionaryUnit)

// Sets a flag to show that a unit has a leaf as a child.
func dictSetHasLeaf(base *DictionaryUnit) {
	*base |= hasLeafBit
}

// Sets a value to a leaf unit.
func dictSetValue(base *DictionaryUnit, value valueType) {
	*base = baseType(value) | isLeafBit
}

// Sets a label to a non-leaf unit.
func dictSetLabel(base *DictionaryUnit, label ucharType) {
	*base = (*base &^ 0xff) | baseType(label)
}

// Sets an offset to a non-leaf unit.
func dictSetOffset(base *DictionaryUnit, offset baseType) bool {
	if offset >= offsetMax<<8 {
		return false
	}

	*base &= isLeafBit | hasLeafBit | 0xff
	if offset < offsetMax {
		*base |= offset << 10
	} else {
		*base |= (offset << 2) | extensionBit
	}
	return true
}

// Checks if a unit has a leaf as a child or not.
func dictHasLeaf(base DictionaryUnit) bool {
	return base&hasLeafBit != 0
}

// Checks if a unit corresponds to a leaf or not.
func dictValue(base DictionaryUnit) valueType {
	return valueType(base &^ isLeafBit)
}

// Reads a label with a leaf flag from a non-leaf unit.
func dictLabel(base DictionaryUnit) ucharType {
	return ucharType(base & (isLeafBit | 0xff))
}

// Reads an offset to child units from a non-leaf unit.
func dictOffset(base DictionaryUnit) baseType {
	return (base >> 10) << ((base & extensionBit) >> 6)
}
