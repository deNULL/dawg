package main

type DictionaryExtraUnit struct {
	loValues baseType
	hiValues baseType
}

func (eu *DictionaryExtraUnit) clear() {
	eu.loValues = 0
	eu.hiValues = 0
}

// Sets if a unit is fixed or not.
func (eu *DictionaryExtraUnit) setIsFixed() {
	eu.loValues |= 1
}

// Sets an index of the next unused unit.
func (eu *DictionaryExtraUnit) setNext(next baseType) {
	eu.loValues = (eu.loValues & 1) | (next << 1)
}

// Sets if an index is used as an offset or not.
func (eu *DictionaryExtraUnit) setIsUsed() {
	eu.hiValues |= 1
}

// Sets an index of the previous unused unit.
func (eu *DictionaryExtraUnit) setPrev(prev baseType) {
	eu.hiValues = (eu.hiValues & 1) | (prev << 1)
}

// Reads if a unit is fixed or not.
func (eu *DictionaryExtraUnit) isFixed() bool {
	return (eu.loValues & 1) == 1
}

// Reads an index of the next unused unit.
func (eu *DictionaryExtraUnit) next() baseType {
	return eu.loValues >> 1
}

// Reads if an index is used as an offset or not.
func (eu *DictionaryExtraUnit) isUsed() bool {
	return (eu.hiValues & 1) == 1
}

// Reads an index of the previous unused unit.
func (eu *DictionaryExtraUnit) prev() baseType {
	return eu.hiValues >> 1
}
