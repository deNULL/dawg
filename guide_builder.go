package dawg

type GuideBuilder struct {
	dawg  *Dawg
	dict  *Dictionary
	guide *Guide

	units        []GuideUnit
	isFixedTable []ucharType // why not use bit pool instead?
}

func NewGuideBuilder(dawg *Dawg, dict *Dictionary, guide *Guide) *GuideBuilder {
	return &GuideBuilder{
		dawg:  dawg,
		dict:  dict,
		guide: guide,
	}
}

func BuildGuide(dawg *Dawg, dict *Dictionary) *Guide {
	builder := NewGuideBuilder(dawg, dict, &Guide{})
	if !builder.Build() {
		return nil
	}
	return builder.guide
}

func (gb *GuideBuilder) Build() bool {
	// Initializes units and flags.
	gb.units = make([]GuideUnit, gb.dict.size)
	gb.isFixedTable = make([]ucharType, gb.dict.size/8)

	if gb.dawg.Size() <= 1 {
		return true
	}

	if !gb.buildIndices(gb.dawg.Root(), gb.dict.Root()) {
		return false
	}

	gb.guide.setUnits(gb.units)
	return true
}

// Builds a guide recursively.
func (gb *GuideBuilder) buildIndices(dawgIndex baseType, dictIndex baseType) bool {
	if gb.isFixed(dictIndex) {
		return true
	}
	gb.setIsFixed(dictIndex)

	// Finds the first non-terminal child.
	var dawgChildIndex baseType = gb.dawg.Child(dawgIndex)
	if gb.dawg.Label(dawgChildIndex) == 0 {
		dawgChildIndex = gb.dawg.Sibling(dawgChildIndex)
		if dawgChildIndex == 0 {
			return true
		}
	}
	gb.units[dictIndex].Child = gb.dawg.Label(dawgChildIndex)

	for {
		var childLabel ucharType = gb.dawg.Label(dawgChildIndex)
		var dictChildIndex baseType = dictIndex
		if !gb.dict.Follow(childLabel, &dictChildIndex) {
			return false
		}

		if !gb.buildIndices(dawgChildIndex, dictChildIndex) {
			return false
		}

		var dawgSiblingIndex baseType = gb.dawg.Sibling(dawgChildIndex)
		var siblingLabel ucharType = gb.dawg.Label(dawgSiblingIndex)
		if dawgSiblingIndex != 0 {
			gb.units[dictChildIndex].Sibling = siblingLabel
		}

		dawgChildIndex = dawgSiblingIndex
		if dawgChildIndex == 0 {
			return true
		}
	}
}

func (gb *GuideBuilder) setIsFixed(index baseType) {
	gb.isFixedTable[index/8] |= 1 << (index % 8)
}

func (gb *GuideBuilder) isFixed(index baseType) bool {
	return gb.isFixedTable[index/8]&(1<<(index%8)) != 0
}
