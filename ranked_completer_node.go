package dawg

type RankedCompleterNode struct {
	dictIndex     baseType
	prevNodeIndex baseType
	label         ucharType
	isQueued      bool
	hasTerminal   bool
}
