package dawg

type RankedGuideLink struct {
	label ucharType
	value valueType
}

func makeRankedGuideLinkCmp(valuesCmp valueComparatorFunc) func(lhs *RankedGuideLink, rhs *RankedGuideLink) bool {
	return func(lhs *RankedGuideLink, rhs *RankedGuideLink) bool {
		if lhs.value != rhs.value {
			return valuesCmp(rhs.value, lhs.value)
		}
		return lhs.label < rhs.label
	}
}
