package main

import "sort"

type RankedGuideBuilder struct {
	dawg  *Dawg
	dict  *Dictionary
	guide *RankedGuide

	units        []RankedGuideUnit
	links        []RankedGuideLink
	isFixedTable []ucharType
}

func BuildRankedGuideCmp(dawg *Dawg, dict *Dictionary, valuesCmp valueComparatorFunc) *RankedGuide {
	var builder *RankedGuideBuilder = &RankedGuideBuilder{
		dawg:  dawg,
		dict:  dict,
		guide: &RankedGuide{},
	}
	if !builder.Build(valuesCmp) {
		return nil
	}
	return builder.guide
}
func BuildRankedGuide(dawg *Dawg, dict *Dictionary) *RankedGuide {
	return BuildRankedGuideCmp(dawg, dict, func(lhs valueType, rhs valueType) bool {
		return lhs < rhs
	})
}

func (rgb *RankedGuideBuilder) Build(valuesCmp valueComparatorFunc) bool {
	// Initializes units and flags.
	rgb.units = make([]RankedGuideUnit, rgb.dict.size)
	rgb.isFixedTable = make([]ucharType, rgb.dict.size/8)

	if rgb.dawg.size() <= 1 {
		return true
	}

	var maxValue valueType = -1
	if !rgb.buildIndices(rgb.dawg.Root(), rgb.dict.Root(), &maxValue, valuesCmp) {
		return false
	}

	rgb.guide.setUnits(rgb.units)
	return true
}

// Builds a guide recursively.
func (rgb *RankedGuideBuilder) buildIndices(dawgIndex baseType, dictIndex baseType, maxValue *valueType, valuesCmp valueComparatorFunc) bool {
	if rgb.isFixed(dictIndex) {
		return rgb.findMaxValue(dictIndex, maxValue)
	}
	rgb.setIsFixed(dictIndex)

	var initialNumLinks sizeType = len(rgb.links)

	// Enumerates links to the next states.
	if !rgb.enumerateLinks(dawgIndex, dictIndex, valuesCmp) {
		return false
	}

	linksCmp := makeRankedGuideLinkCmp(valuesCmp)
	sort.SliceStable(rgb.links[initialNumLinks:], func(i int, j int) bool {
		return linksCmp(&rgb.links[initialNumLinks+i], &rgb.links[initialNumLinks+j])
	})

	// Reflects links into units.
	if !rgb.turnLinksToUnits(dictIndex, initialNumLinks) {
		return false
	}

	*maxValue = rgb.links[initialNumLinks].value
	rgb.links = rgb.links[:initialNumLinks]

	return true
}

// Finds the maximum value by using fixed units.
func (rgb *RankedGuideBuilder) findMaxValue(dictIndex baseType, maxValue *valueType) bool {
	for rgb.units[dictIndex].Child != 0 {
		var childLabel ucharType = rgb.units[dictIndex].Child
		if !rgb.dict.Follow(childLabel, &dictIndex) {
			return false
		}
	}
	if !rgb.dict.HasValue(dictIndex) {
		return false
	}
	*maxValue = rgb.dict.Value(dictIndex)
	return true
}

// Enumerates links to the next states.
func (rgb *RankedGuideBuilder) enumerateLinks(dawgIndex baseType, dictIndex baseType, valuesCmp valueComparatorFunc) bool {
	for dawgChildIndex := rgb.dawg.Child(dawgIndex); dawgChildIndex != 0; dawgChildIndex = rgb.dawg.Sibling(dawgChildIndex) {
		var value valueType = -1
		var childLabel ucharType = rgb.dawg.Label(dawgChildIndex)
		if childLabel == 0 {
			if !rgb.dict.HasValue(dictIndex) {
				return false
			}
			value = rgb.dict.Value(dictIndex)
		} else {
			var dictChildIndex = dictIndex
			if !rgb.dict.Follow(childLabel, &dictChildIndex) {
				return false
			}

			if !rgb.buildIndices(dawgChildIndex, dictChildIndex, &value, valuesCmp) {
				return false
			}
		}
		rgb.links = append(rgb.links, RankedGuideLink{
			label: childLabel,
			value: value,
		})
	}

	return true
}

// Modifies units.
func (rgb *RankedGuideBuilder) turnLinksToUnits(dictIndex baseType, linksBegin sizeType) bool {
	// The first child.
	var firstLabel ucharType = rgb.links[linksBegin].label
	rgb.units[dictIndex].Child = firstLabel
	var dictChildIndex baseType = rgb.followWithoutCheck(dictIndex, firstLabel)

	// Other children.
	for i := linksBegin + 1; i < len(rgb.links); i++ {
		var siblingLabel ucharType = rgb.links[i].label

		var dictSiblingIndex = rgb.followWithoutCheck(dictIndex, siblingLabel)
		rgb.units[dictChildIndex].Sibling = siblingLabel
		dictChildIndex = dictSiblingIndex
	}

	return true
}

// Follows a transition without any check.
func (rgb *RankedGuideBuilder) followWithoutCheck(index baseType, label ucharType) baseType {
	return index ^ dictOffset(rgb.dict.units[index]) ^ baseType(label)
}

func (rgb *RankedGuideBuilder) setIsFixed(index baseType) {
	rgb.isFixedTable[index/8] |= 1 << (index % 8)
}

func (rgb *RankedGuideBuilder) isFixed(index baseType) bool {
	return rgb.isFixedTable[index/8]&(1<<(index%8)) != 0
}
