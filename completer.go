package main

type SomeCompleter interface {
	Start(baseType)
	Next() bool
	Key() string
	Value() valueType
	Length() sizeType
}
type Completer struct {
	dict  *Dictionary
	guide *Guide

	path       []ucharType
	indexStack []baseType
	lastIndex  baseType
}

func NewCompleter(dict *Dictionary, guide *Guide) *Completer {
	return &Completer{
		dict:  dict,
		guide: guide,
	}
}

// These member functions are available only when Next() returns true.
func (c *Completer) Length() sizeType {
	return len(c.path) - 1
}
func (c *Completer) Key() string {
	return string(c.path)
}
func (c *Completer) Value() valueType {
	return c.dict.Value(c.lastIndex)
}

// Starts completing keys from given index and prefix.
func (c *Completer) Start(index baseType) {
	c.StartStringLen(index, "", 0)
}
func (c *Completer) StartString(index baseType, prefix string) {
	c.StartStringLen(index, prefix, len(prefix))
}
func (c *Completer) StartStringLen(index baseType, prefix string, length sizeType) {
	c.path = make([]ucharType, length+1)
	for i := 0; i < length; i++ {
		c.path[i] = prefix[i]
	}
	c.path[length] = 0

	c.indexStack = c.indexStack[:0]
	if c.guide.size != 0 {
		c.indexStack = append(c.indexStack, index)
		c.lastIndex = c.dict.Root()
	}
}

// Gets the next key.
func (c *Completer) Next() bool {
	if len(c.indexStack) == 0 {
		return false
	}
	var index baseType = c.indexStack[len(c.indexStack)-1]

	if c.lastIndex != c.dict.Root() {
		var childLabel ucharType = c.guide.Child(index)
		if childLabel != 0 {
			// Follows a transition to the first child.
			if !c.Follow(childLabel, &index) {
				return false
			}
		} else {
			for {
				var siblingLabel ucharType = c.guide.Sibling(index)

				// Moves to the previous node
				if len(c.path) > 1 {
					c.path = c.path[:len(c.path)-1]
					c.path[len(c.path)-1] = 0
				}
				c.indexStack = c.indexStack[:len(c.indexStack)-1]
				if len(c.indexStack) == 0 {
					return false
				}

				index = c.indexStack[len(c.indexStack)-1]
				if siblingLabel != 0 {
					// Follows a transition to the next sibling.
					if !c.Follow(siblingLabel, &index) {
						return false
					}
					break
				}
			}
		}
	}

	// Finds a terminal
	return c.FindTerminal(index)
}

// Follows a transition.
func (c *Completer) Follow(label ucharType, index *baseType) bool {
	if !c.dict.Follow(label, index) {
		return false
	}

	c.path[len(c.path)-1] = label
	c.path = append(c.path, 0)
	c.indexStack = append(c.indexStack, *index)
	return true
}

// Finds a terminal
func (c *Completer) FindTerminal(index baseType) bool {
	for !c.dict.HasValue(index) {
		var label ucharType = c.guide.Child(index)
		if !c.dict.Follow(label, &index) {
			return false
		}

		c.path[len(c.path)-1] = label
		c.path = append(c.path, 0)
		c.indexStack = append(c.indexStack, index)
	}

	c.lastIndex = index
	return true
}
