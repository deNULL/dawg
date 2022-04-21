package dawg

type IndexBuilder struct {
	dict  *Dictionary
	guide SomeGuide
	index *Index
}

func NewIndexBuilder(dict *Dictionary, guide SomeGuide, index *Index) *IndexBuilder {
	return &IndexBuilder{
		dict:  dict,
		guide: guide,
		index: index,
	}
}

func BuildIndex(dict *Dictionary, guide SomeGuide) *Index {
	builder := NewIndexBuilder(dict, guide, &Index{})
	if !builder.Build() {
		return nil
	}
	return builder.index
}

func (ib *IndexBuilder) Build() bool {
	ib.index.units = make([]IndexUnit, len(ib.dict.units))
	return ib.buildIndices(ib.guide.Root())
}

func (ib *IndexBuilder) buildIndices(index baseType) bool {
	if ib.dict.HasValue(index) {
		ib.index.units[index]++
	}

	var child ucharType = ib.guide.Child(index)
	for child != 0 {
		var childIndex baseType = index
		if !ib.dict.Follow(child, &childIndex) {
			return false
		}
		if ib.index.units[childIndex] == 0 {
			if !ib.buildIndices(childIndex) {
				return false
			}
		}
		ib.index.units[index] += ib.index.units[childIndex]
		child = ib.guide.Sibling(childIndex)
	}

	return true
}
