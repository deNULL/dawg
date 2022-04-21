package dawg

type Indexer struct {
	dict  *Dictionary
	guide SomeGuide
	index *Index
}

func NewIndexer(dict *Dictionary, guide SomeGuide, index *Index) *Indexer {
	return &Indexer{
		dict:  dict,
		guide: guide,
		index: index,
	}
}

const NotFound baseType = 0xffffffff
const Failed baseType = 0xfffffffe

func (idx *Indexer) TotalCount() baseType {
	return idx.index.TotalCount()
}

func (idx *Indexer) StringToIndex(s string) baseType {
	return idx.BytesToIndex(([]ucharType)(s))
}

func (idx *Indexer) IndexToString(i baseType) string {
	return string(idx.IndexToBytes(i))
}

func (idx *Indexer) BytesToIndex(bytes []ucharType) baseType {
	var index baseType = idx.dict.Root()
	var result baseType = 0
	for i := 0; i < len(bytes); i++ {
		if idx.dict.HasValue(index) {
			result++
		}

		var childLabel ucharType = idx.guide.Child(index)
		for childLabel < bytes[i] && childLabel != 0 {
			var childIndex baseType = index
			if !idx.dict.Follow(childLabel, &childIndex) {
				return Failed
			}
			result += idx.index.ChildCount(childIndex)
			childLabel = idx.guide.Sibling(childIndex)
		}

		if childLabel == 0 {
			return NotFound
		}

		var childIndex baseType = index
		if !idx.dict.Follow(bytes[i], &childIndex) {
			return NotFound
		}
		index = childIndex
	}

	if !idx.dict.HasValue(index) {
		return NotFound
	}

	return result
}

func (idx *Indexer) IndexToBytes(i baseType) []ucharType {
	var index baseType = idx.dict.Root()
	var cur baseType = 0
	var buf []ucharType = make([]ucharType, 0)
	for cur <= i {
		if idx.dict.HasValue(index) {
			if cur == i {
				return buf
			}
			cur++
		}

		var childLabel ucharType = idx.guide.Child(index)
		for childLabel != 0 {
			var childIndex baseType = index
			if !idx.dict.Follow(childLabel, &childIndex) {
				return nil
			}
			var count baseType = idx.index.ChildCount(childIndex)
			if i < cur+count {
				buf = append(buf, childLabel)
				index = childIndex
				break
			} else {
				cur += count
				childLabel = idx.guide.Sibling(childIndex)
			}
		}

		if childLabel == 0 {
			return nil
		}
	}

	return nil
}
