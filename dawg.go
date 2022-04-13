package main

type Dawg struct {
	basePool               []BaseUnit
	labelPool              []ucharType
	flagPool               *BitPool
	numOfStates            sizeType
	numOfMergedTransitions sizeType
	numOfMergedStates      sizeType
	numOfMergingStates     sizeType
}

func NewDawg() *Dawg {
	return &Dawg{
		basePool:               []BaseUnit{},
		labelPool:              []ucharType{},
		flagPool:               NewBitPool(),
		numOfStates:            0,
		numOfMergedTransitions: 0,
		numOfMergedStates:      0,
		numOfMergingStates:     0,
	}
}

// The root index.
func (dawg *Dawg) Root() baseType {
	return 0
}

// Number of units.
func (dawg *Dawg) size() sizeType {
	return len(dawg.basePool)
}

// Number of transitions.
func (dawg *Dawg) numOfTransitions() sizeType {
	return len(dawg.basePool) - 1
}

// Reads values.
func (dawg *Dawg) Child(index baseType) baseType {
	return dawg.basePool[index].child()
}
func (dawg *Dawg) Sibling(index baseType) baseType {
	if dawg.basePool[index].hasSibling() {
		return index + 1
	}
	return 0
}
func (dawg *Dawg) Value(index baseType) valueType {
	return dawg.basePool[index].value()
}

func (dawg *Dawg) IsLeaf(index baseType) bool {
	return dawg.Label(index) == 0
}
func (dawg *Dawg) Label(index baseType) ucharType {
	return dawg.labelPool[index]
}
func (dawg *Dawg) IsMerging(index baseType) bool {
	return dawg.flagPool.get(index)
}

// Clears object pools.
func (dawg *Dawg) Clear() {
	dawg.basePool = dawg.basePool[:0]
	dawg.labelPool = dawg.labelPool[:0]
	dawg.flagPool.clear()
	dawg.numOfStates = 0
	dawg.numOfMergedStates = 0
}
