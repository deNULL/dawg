package main

const defaultInitialHashTableSize = 1 << 8

type DawgBuilder struct {
	initialHashTableSize   sizeType
	basePool               []BaseUnit
	labelPool              []ucharType
	flagPool               *BitPool
	unitPool               []DawgUnit
	hashTable              []baseType
	unfixedUnits           []baseType
	unusedUnits            []baseType
	numOfStates            sizeType
	numOfMergedTransitions sizeType
	numOfMergingStates     sizeType
}

func NewDawgBuilder() *DawgBuilder {
	return NewDawgBuilderWithSize(defaultInitialHashTableSize)
}
func NewDawgBuilderWithSize(initialHashTableSize sizeType) *DawgBuilder {
	return &DawgBuilder{
		initialHashTableSize:   initialHashTableSize,
		basePool:               []BaseUnit{},
		labelPool:              []ucharType{},
		flagPool:               NewBitPool(),
		unitPool:               []DawgUnit{},
		hashTable:              []baseType{},
		unfixedUnits:           []baseType{},
		unusedUnits:            []baseType{},
		numOfStates:            1,
		numOfMergedTransitions: 0,
		numOfMergingStates:     0,
	}
}

func (db *DawgBuilder) numOfTransitions() sizeType {
	return len(db.basePool) - 1
}

func (db *DawgBuilder) numOfMergedStates() sizeType {
	return db.numOfTransitions() + db.numOfMergedTransitions + 1 - db.numOfStates
}

// Gets a unit from an object pool.
func (db *DawgBuilder) allocateUnit() baseType {
	var index baseType = 0
	if len(db.unusedUnits) == 0 {
		index = baseType(len(db.unitPool))
		db.unitPool = append(db.unitPool, DawgUnit{})
	} else {
		index = db.unusedUnits[len(db.unusedUnits)-1]
		db.unusedUnits = db.unusedUnits[:len(db.unusedUnits)-1]
		db.unitPool[index].clear()
	}
	return index
}

// Returns a unit to an object pool.
func (db *DawgBuilder) freeUnit(index baseType) {
	db.unusedUnits = append(db.unusedUnits, index)
}

// Gets a transition from object pools.
func (db *DawgBuilder) allocateTransition() baseType {
	db.flagPool.allocate()
	db.basePool = append(db.basePool, 0)
	db.labelPool = append(db.labelPool, 0)
	return baseType(len(db.labelPool) - 1)
}

func (db *DawgBuilder) init() {
	db.hashTable = make([]baseType, db.initialHashTableSize)
	db.allocateUnit()
	db.allocateTransition()
	db.unitPool[0].label = 0xff
	db.unfixedUnits = append(db.unfixedUnits, 0)
}

// 32-bit mix function.
// http://www.concentric.net/~Ttwang/tech/inthash.htm
func intHash(key baseType) baseType {
	key = ^key + (key << 15) // key = (key << 15) - key - 1;
	key = key ^ (key >> 12)
	key = key + (key << 2)
	key = key ^ (key >> 4)
	key = key * 2057 // key = (key + (key << 3)) + (key << 11);
	key = key ^ (key >> 16)
	return key
}

// Calculates a hash value from a transition.
func (db *DawgBuilder) hashTransition(index baseType) baseType {
	var hashValue baseType = 0
	for index != 0 {
		var base baseType = baseType(db.basePool[index])
		var label ucharType = db.labelPool[index]
		hashValue ^= intHash((baseType(label) << 24) ^ base)

		if !db.basePool[index].hasSibling() {
			break
		}
		index += 1
	}
	return hashValue
}

// Calculates a hash value from a unit.
func (db *DawgBuilder) hashUnit(index baseType) baseType {
	var hashValue baseType = 0
	for index != 0 {
		var base baseType = db.unitPool[index].base()
		var label ucharType = db.unitPool[index].label
		hashValue ^= intHash((baseType(label) << 24) ^ base)
		index = db.unitPool[index].sibling
	}
	return hashValue
}

// Compares a unit and a transition.
func (db *DawgBuilder) areEqual(unitIndex baseType, transitionIndex baseType) bool {
	// Compares the numbers of transitions.
	for i := db.unitPool[unitIndex].sibling; i != 0; i = db.unitPool[i].sibling {
		if !db.basePool[transitionIndex].hasSibling() {
			return false
		}
		transitionIndex += 1
	}
	if db.basePool[transitionIndex].hasSibling() {
		return false
	}

	// Compares out-transitions.
	for i := unitIndex; i != 0; i = db.unitPool[i].sibling {
		if db.unitPool[i].base() != baseType(db.basePool[transitionIndex]) ||
			db.unitPool[i].label != db.labelPool[transitionIndex] {
			return false
		}
		transitionIndex -= 1
	}
	return true
}

// Finds a transition from a hash table.
func (db *DawgBuilder) findTransition(index baseType, hashId *baseType) baseType {
	*hashId = db.hashTransition(index) % baseType(len(db.hashTable))
	for {
		var transitionId baseType = db.hashTable[*hashId]
		if transitionId == 0 {
			break
		}

		// There must not be the same base value.
		*hashId = (*hashId + 1) % baseType(len(db.hashTable))
	}
	return 0
}

// Finds a unit from a hash table.
func (db *DawgBuilder) findUnit(unitIndex baseType, hashId *baseType) baseType {
	*hashId = db.hashUnit(unitIndex) % baseType(len(db.hashTable))
	for {
		var transitionId baseType = db.hashTable[*hashId]
		if transitionId == 0 {
			break
		}

		if db.areEqual(unitIndex, transitionId) {
			return transitionId
		}
		*hashId = (*hashId + 1) % baseType(len(db.hashTable))
	}
	return 0
}

// Expands a hash table.
func (db *DawgBuilder) expandHashTable() {
	var hashTableSize sizeType = len(db.hashTable) << 1
	db.hashTable = make([]baseType, hashTableSize)

	// Builds a new hash table.
	var count baseType = 0
	for i := 1; i < len(db.basePool); i += 1 {
		var index = baseType(i)
		if db.labelPool[index] == 0 || db.basePool[index].isState() {
			var hashId baseType
			db.findTransition(index, &hashId)
			db.hashTable[hashId] = index
			count += 1
		}
	}
}

// Fixes units corresponding to the last inserted key.
// Also, some of units are merged into equivalent transitions.
func (db *DawgBuilder) fixUnits(index baseType) {
	for db.unfixedUnits[len(db.unfixedUnits)-1] != index {
		var unfixedIndex baseType = db.unfixedUnits[len(db.unfixedUnits)-1]
		db.unfixedUnits = db.unfixedUnits[:len(db.unfixedUnits)-1]

		if db.numOfStates >= len(db.hashTable)-len(db.hashTable)>>2 {
			db.expandHashTable()
		}

		var numOfSiblings sizeType = 0
		for i := unfixedIndex; i != 0; i = db.unitPool[i].sibling {
			numOfSiblings += 1
		}

		var hashId baseType
		var matchedIndex baseType = db.findUnit(unfixedIndex, &hashId)
		if matchedIndex != 0 {
			db.numOfMergedTransitions += numOfSiblings

			// Records a merging state.
			if !db.flagPool.get(matchedIndex) {
				db.numOfMergingStates += 1
				db.flagPool.set(matchedIndex, true)
			}
		} else {
			// Fixes units into pairs of base values and labels.
			var transitionIndex baseType = 0
			for i := 0; i < numOfSiblings; i += 1 {
				transitionIndex = db.allocateTransition()
			}
			for i := unfixedIndex; i != 0; i = db.unitPool[i].sibling {
				db.basePool[transitionIndex] = BaseUnit(db.unitPool[i].base())
				db.labelPool[transitionIndex] = db.unitPool[i].label
				transitionIndex -= 1
			}
			matchedIndex = transitionIndex + 1
			db.hashTable[hashId] = matchedIndex
			db.numOfStates += 1
		}

		// Deletes fixed units.
		var next baseType
		for current := unfixedIndex; current != 0; current = next {
			next = db.unitPool[current].sibling
			db.freeUnit(current)
		}

		db.unitPool[db.unfixedUnits[len(db.unfixedUnits)-1]].child = matchedIndex
	}
	db.unfixedUnits = db.unfixedUnits[:len(db.unfixedUnits)-1]
}

func (db *DawgBuilder) clear() {
	db.basePool = db.basePool[:0]
	db.labelPool = db.labelPool[:0]
	db.flagPool.clear()
	db.unitPool = db.unitPool[:0]

	db.hashTable = []baseType{0}
	db.unfixedUnits = db.unfixedUnits[:0]
	db.unusedUnits = db.unusedUnits[:0]

	db.numOfStates = 1
	db.numOfMergedTransitions = 0
	db.numOfMergingStates = 0
}

func (db *DawgBuilder) InsertKeyValue(key []ucharType, length sizeType, value valueType) bool {
	// Initializes a builder if not initialized.
	if len(db.hashTable) == 0 {
		db.init()
	}

	var index baseType = 0
	var keyPos sizeType = 0

	// Finds a separate unit.
	for keyPos <= length {
		var childIndex baseType = db.unitPool[index].child
		if childIndex == 0 {
			break
		}

		var keyLabel ucharType = 0
		if keyPos < length {
			keyLabel = key[keyPos]
		}
		var unitLabel ucharType = db.unitPool[childIndex].label

		// Checks the order of keys.
		if keyLabel < unitLabel {
			return false
		} else if keyLabel > unitLabel {
			db.unitPool[childIndex].hasSibling = true
			db.fixUnits(childIndex)
			break
		}

		index = childIndex
		keyPos += 1
	}

	// Adds new units.
	for keyPos <= length {
		var keyLabel ucharType = 0
		if keyPos < length {
			keyLabel = key[keyPos]
		}
		var childIndex baseType = db.allocateUnit()

		if db.unitPool[index].child == 0 {
			db.unitPool[childIndex].isState = true
		}
		db.unitPool[childIndex].sibling = db.unitPool[index].child
		db.unitPool[childIndex].label = keyLabel
		db.unitPool[index].child = childIndex
		db.unfixedUnits = append(db.unfixedUnits, childIndex)

		index = childIndex
		keyPos += 1
	}
	db.unitPool[index].setValue(value)
	return true
}

func (db *DawgBuilder) InsertStringValue(key string, value valueType) bool {
	return db.InsertKeyValue([]ucharType(key), len(key), value)
}

func (db *DawgBuilder) InsertString(key string) bool {
	return db.InsertKeyValue([]ucharType(key), len(key), 0)
}

// Finishes building a dawg.
func (db *DawgBuilder) Finish(dawg *Dawg) {
	// Initializes a builder if not initialized.
	if len(db.hashTable) == 0 {
		db.init()
	}

	db.fixUnits(0)
	db.basePool[0] = BaseUnit(db.unitPool[0].base())
	db.labelPool[0] = db.unitPool[0].label

	dawg.numOfStates = db.numOfStates
	dawg.numOfMergedTransitions = db.numOfMergedTransitions
	dawg.numOfMergedStates = db.numOfMergedStates()
	dawg.numOfMergingStates = db.numOfMergingStates

	dawg.basePool, db.basePool = db.basePool, dawg.basePool
	dawg.labelPool, db.labelPool = db.labelPool, dawg.labelPool
	dawg.flagPool, db.flagPool = db.flagPool, dawg.flagPool

	db.clear()
}
