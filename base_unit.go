package dawg

type BaseUnit baseType

func (unit *BaseUnit) child() baseType {
	return baseType(*unit >> 2)
}

func (unit *BaseUnit) hasSibling() bool {
	return *unit&1 != 0
}

func (unit *BaseUnit) isState() bool {
	return *unit&2 != 0
}

func (unit *BaseUnit) value() valueType {
	return valueType(*unit >> 1)
}
