package main

type linkPair struct {
	index  baseType
	offset baseType
}

type LinkTable struct {
	hashTable []linkPair
}

func NewLinkTable(tableSize sizeType) *LinkTable {
	return &LinkTable{
		hashTable: make([]linkPair, tableSize),
	}
}

// Finds an Id from an upper table.
func (lt *LinkTable) findId(index baseType) baseType {
	var hashId baseType = intHash(index) % baseType(len(lt.hashTable))
	for lt.hashTable[hashId].index != 0 {
		if index == lt.hashTable[hashId].index {
			return hashId
		}
		hashId = (hashId + 1) % baseType(len(lt.hashTable))
	}
	return hashId
}

// Inserts an index with its offset.
func (lt *LinkTable) Insert(index baseType, offset baseType) {
	var hashId = lt.findId(index)
	lt.hashTable[hashId].index = index
	lt.hashTable[hashId].offset = offset
}

// Finds an offset that corresponds to a given index.
func (lt *LinkTable) Find(index baseType) baseType {
	var hashId = lt.findId(index)
	return lt.hashTable[hashId].offset
}
