package main

// Number of units in a block.
const blockSize = 256

// Number of blocks kept unfixed.
const numOfUnfixedBlocks = 16

// Number of units kept unfixed.
const unfixedSize = blockSize * numOfUnfixedBlocks

// Masks for offsets.
const upperMask = ^(offsetMax - 1)
const lowerMask = 0xFF

type DictionaryBuilder struct {
	dawg *Dawg
	dict *Dictionary

	units            []DictionaryUnit
	extras           [][]DictionaryExtraUnit
	labels           []ucharType
	linkTable        *LinkTable
	unfixedIndex     baseType
	numOfUnusedUnits baseType
}

func NewDictionaryBuilder(dawg *Dawg, dict *Dictionary) *DictionaryBuilder {
	return &DictionaryBuilder{
		dawg: dawg,
		dict: dict,
	}
}

func (dawg *Dawg) Build() *Dictionary {
	builder := NewDictionaryBuilder(dawg, NewDictionary())
	if !builder.BuildDictionary() {
		return nil
	}
	return builder.dict
}

func (dawg *Dawg) BuildWithUnused(numOfUnusedUnits *baseType) *Dictionary {
	builder := NewDictionaryBuilder(dawg, NewDictionary())
	if !builder.BuildDictionary() {
		return nil
	}
	*numOfUnusedUnits = builder.numOfUnusedUnits
	return builder.dict
}

func (db *DictionaryBuilder) numOfUnits() baseType {
	return baseType(len(db.units))
}

func (db *DictionaryBuilder) numOfBlocks() baseType {
	return baseType(len(db.extras))
}

func (db *DictionaryBuilder) extra(index baseType) *DictionaryExtraUnit {
	return &db.extras[index/blockSize][index%blockSize]
}

// Builds a dictionary from a list-form dawg.
func (db *DictionaryBuilder) BuildDictionary() bool {
	db.linkTable = NewLinkTable(db.dawg.numOfMergingStates + (db.dawg.numOfMergingStates >> 1))
	db.reserveUnit(0)
	db.extra(0).setIsUsed()
	dictSetOffset(&db.units[0], 1)
	dictSetLabel(&db.units[0], 0)

	if db.dawg.size() > 1 {
		if !db.buildDictionaryIndices(db.dawg.Root(), 0) {
			return false
		}
	}

	db.fixAllBlocks()
	db.dict.setUnits(db.units)
	return true
}

// Builds a dictionary from a dawg.
func (db *DictionaryBuilder) buildDictionaryIndices(dawgIndex baseType, dictIndex baseType) bool {
	if db.dawg.IsLeaf(dawgIndex) {
		return true
	}

	// Uses an existing offset if available.
	var dawgChildIndex baseType = db.dawg.Child(dawgIndex)
	if db.dawg.IsMerging(dawgChildIndex) {
		var offset baseType = db.linkTable.Find(dawgChildIndex)
		if offset != 0 {
			offset ^= dictIndex
			if (offset&upperMask == 0) || (offset&lowerMask == 0) {
				if db.dawg.IsLeaf(dawgChildIndex) {
					dictSetHasLeaf(&db.units[dictIndex])
				}
				dictSetOffset(&db.units[dictIndex], offset)
				return true
			}
		}
	}

	// Finds a good offset and arranges child nodes.
	var offset baseType = db.arrangeChildNodes(dawgIndex, dictIndex)
	if offset == 0 {
		return false
	}

	if db.dawg.IsMerging(dawgChildIndex) {
		db.linkTable.Insert(dawgChildIndex, offset)
	}

	// Builds a double-array in depth-first order.
	for {
		var dictChildIndex baseType = offset ^ baseType(db.dawg.Label(dawgChildIndex))
		if !db.buildDictionaryIndices(dawgChildIndex, dictChildIndex) {
			return false
		}
		dawgChildIndex = db.dawg.Sibling(dawgChildIndex)
		if dawgChildIndex == 0 {
			break
		}
	}
	return true
}

// Arranges child nodes.
func (db *DictionaryBuilder) arrangeChildNodes(dawgIndex baseType, dictIndex baseType) baseType {
	db.labels = db.labels[:0]

	var dawgChildIndex baseType = db.dawg.Child(dawgIndex)
	for dawgChildIndex != 0 {
		db.labels = append(db.labels, db.dawg.Label(dawgChildIndex))
		dawgChildIndex = db.dawg.Sibling(dawgChildIndex)
	}

	// Finds a good offset.
	var offset baseType = db.findGoodOffset(dictIndex)
	if !dictSetOffset(&db.units[dictIndex], dictIndex^offset) {
		return 0
	}

	dawgChildIndex = db.dawg.Child(dawgIndex)
	for i := 0; i < len(db.labels); i++ {
		var dictChildIndex baseType = offset ^ baseType(db.labels[i])
		db.reserveUnit(dictChildIndex)

		if db.dawg.IsLeaf(dawgChildIndex) {
			dictSetHasLeaf(&db.units[dictIndex])
			dictSetValue(&db.units[dictChildIndex], db.dawg.Value(dawgChildIndex))
		} else {
			dictSetLabel(&db.units[dictChildIndex], db.labels[i])
		}

		dawgChildIndex = db.dawg.Sibling(dawgChildIndex)
	}
	db.extra(offset).setIsUsed()

	return offset
}

// Finds a good offset.
func (db *DictionaryBuilder) findGoodOffset(index baseType) baseType {
	if db.unfixedIndex >= db.numOfUnits() {
		return db.numOfUnits() | (index & 0xff)
	}

	// Scans unused units to find a good offset.
	var unfixedIndex baseType = db.unfixedIndex
	for {
		var offset baseType = unfixedIndex ^ baseType(db.labels[0])
		if db.isGoodOffset(index, offset) {
			return offset
		}
		unfixedIndex = db.extra(unfixedIndex).next()
		if unfixedIndex == db.unfixedIndex {
			break
		}
	}

	return db.numOfUnits() | (index & 0xff)
}

// Checks if a given offset is valid or not.
func (db *DictionaryBuilder) isGoodOffset(index baseType, offset baseType) bool {
	if db.extra(offset).isUsed() {
		return false
	}

	var relativeOffset baseType = index ^ offset
	if (relativeOffset&lowerMask != 0) && (relativeOffset&upperMask != 0) {
		return false
	}

	// Finds a collision
	for i := 1; i < len(db.labels); i++ {
		if db.extra(offset ^ baseType(db.labels[i])).isFixed() {
			return false
		}
	}

	return true
}

// Reserves an unused unit.
func (db *DictionaryBuilder) reserveUnit(index baseType) {
	if index >= db.numOfUnits() {
		db.expandDictionary()
	}

	// Removes an unused unit from a circular linked list.
	if index == db.unfixedIndex {
		db.unfixedIndex = db.extra(index).next()
		if db.unfixedIndex == index {
			db.unfixedIndex = db.numOfUnits()
		}
	}
	db.extra(db.extra(index).prev()).setNext(db.extra(index).next())
	db.extra(db.extra(index).next()).setPrev(db.extra(index).prev())
	db.extra(index).setIsFixed()
}

// Expands a dictionary.
func (db *DictionaryBuilder) expandDictionary() {
	var srcNumOfUnits baseType = db.numOfUnits()
	var srcNumOfBlocks baseType = db.numOfBlocks()

	var destNumOfUnits baseType = srcNumOfUnits + blockSize
	var destNumOfBlocks baseType = srcNumOfBlocks + 1

	// Fixes an old block
	if destNumOfBlocks > numOfUnfixedBlocks {
		db.fixBlock(srcNumOfBlocks - numOfUnfixedBlocks)
	}

	db.units = append(db.units, make([]DictionaryUnit, blockSize)...)
	db.extras = append(db.extras, nil)

	// Allocates memory to a new block.
	if destNumOfBlocks > numOfUnfixedBlocks {
		var blockId baseType = srcNumOfBlocks - numOfUnfixedBlocks
		db.extras[blockId], db.extras[len(db.extras)-1] = db.extras[len(db.extras)-1], db.extras[blockId]
		for i := srcNumOfUnits; i < destNumOfUnits; i++ {
			db.extra(i).clear()
		}
	} else {
		db.extras[len(db.extras)-1] = make([]DictionaryExtraUnit, blockSize)
	}

	// Creates a circular linked list for a new block.
	for i := srcNumOfUnits + 1; i < destNumOfUnits; i++ {
		db.extra(i - 1).setNext(i)
		db.extra(i).setPrev(i - 1)
	}

	db.extra(srcNumOfUnits).setPrev(destNumOfUnits - 1)
	db.extra(destNumOfUnits - 1).setNext(srcNumOfUnits)

	// Merges 2 circular linked lists.
	db.extra(srcNumOfUnits).setPrev(db.extra(db.unfixedIndex).prev())
	db.extra(destNumOfUnits - 1).setNext(db.unfixedIndex)

	db.extra(db.extra(db.unfixedIndex).prev()).setNext(srcNumOfUnits)
	db.extra(db.unfixedIndex).setPrev(destNumOfUnits - 1)
}

// Fixes all blocks to avoid invalid transitions.
func (db *DictionaryBuilder) fixAllBlocks() {
	var begin baseType = 0
	if db.numOfBlocks() > numOfUnfixedBlocks {
		begin = db.numOfBlocks() - numOfUnfixedBlocks
	}
	var end baseType = db.numOfBlocks()

	for blockId := begin; blockId != end; blockId++ {
		db.fixBlock(blockId)
	}
}

// Adjusts labels of unused units in a given block.
func (db *DictionaryBuilder) fixBlock(blockId baseType) {
	var begin baseType = blockId * blockSize
	var end baseType = begin + blockSize

	// Finds an unused offset.
	var unusedOffsetForLabel baseType = 0
	for offset := begin; offset != end; offset++ {
		if !db.extra(offset).isUsed() {
			unusedOffsetForLabel = offset
			break
		}
	}

	// Labels of unused units are modified.
	for index := begin; index != end; index++ {
		if !db.extra(index).isFixed() {
			db.reserveUnit(index)
			dictSetLabel(&db.units[index], ucharType(index^unusedOffsetForLabel))
			db.numOfUnusedUnits += 1
		}
	}
}
