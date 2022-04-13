package main

type DawgUnit struct {
	child      baseType
	sibling    baseType
	label      ucharType
	isState    bool
	hasSibling bool
}

// Calculates a base value of a unit.
func (unit *DawgUnit) base() baseType {
	if unit.label == 0 {
		var base = unit.child << 1
		if unit.hasSibling {
			base |= 1
		}
		return base
	}

	var base = unit.child << 2
	if unit.isState {
		base |= 2
	}
	if unit.hasSibling {
		base |= 1
	}
	return base
}

func (unit *DawgUnit) setValue(value valueType) {
	unit.child = baseType(value)
}

func (unit *DawgUnit) value() valueType {
	return valueType(unit.child)
}

func (unit *DawgUnit) clear() {
	unit.child = 0
	unit.sibling = 0
	unit.label = 0
	unit.isState = false
	unit.hasSibling = false
}
