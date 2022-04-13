package main

import (
	"container/heap"
)

type RankedCompleter struct {
	dict  *Dictionary
	guide *RankedGuide

	path         []ucharType
	prefixLength sizeType
	value        valueType

	nodes          []RankedCompleterNode
	nodeQueue      []baseType
	candidateQueue RankedCompleterCandidateQueue
}

func NewRankedCompleter(dict *Dictionary, guide *RankedGuide) *RankedCompleter {
	return &RankedCompleter{
		dict:  dict,
		guide: guide,

		value: -1,
	}
}

func (rc *RankedCompleter) Key() string {
	return string(rc.path)
}

func (rc *RankedCompleter) Value() valueType {
	return rc.value
}

func (rc *RankedCompleter) Length() sizeType {
	return len(rc.path) - 1
}

func (rc *RankedCompleter) Start(index baseType) {
	rc.StartStringLen(index, "", 0)
}
func (rc *RankedCompleter) StartString(index baseType, prefix string) {
	rc.StartStringLen(index, prefix, len(prefix))
}
func (rc *RankedCompleter) StartStringLen(index baseType, prefix string, length sizeType) {
	rc.path = make([]ucharType, length)
	for i := 0; i < length; i++ {
		rc.path[i] = prefix[i]
	}
	rc.prefixLength = length
	rc.value = -1

	rc.nodes = rc.nodes[:0]
	rc.nodeQueue = rc.nodeQueue[:0]
	rc.candidateQueue = rc.candidateQueue[:0]

	if rc.guide.size != 0 {
		rc.createNode(index, 0, 'X')
		rc.enqueueNode(0)
	}
}

func reverseUCharSlice(slice []ucharType) {
	last := len(slice) - 1
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[last-i] = slice[last-i], slice[i]
	}
}

// Gets the next key.
func (rc *RankedCompleter) Next() bool {
	for i := 0; i < len(rc.nodeQueue); i++ {
		var nodeIndex baseType = rc.nodeQueue[i]
		if rc.value != -1 && !rc.findSibling(&nodeIndex) {
			continue
		}
		nodeIndex = rc.findTerminal(nodeIndex)
		rc.enqueueCandidate(nodeIndex)
	}
	rc.nodeQueue = rc.nodeQueue[:0]

	// Returns false if there is no candidate.
	if len(rc.candidateQueue) == 0 {
		return false
	}

	var candidate *RankedCompleterCandidate = rc.candidateQueue[0]

	var nodeIndex baseType = candidate.nodeIndex
	rc.enqueueNode(nodeIndex)
	nodeIndex = rc.nodes[nodeIndex].prevNodeIndex

	rc.path = rc.path[:rc.prefixLength]
	for nodeIndex != 0 {
		rc.path = append(rc.path, rc.nodes[nodeIndex].label)
		rc.enqueueNode(nodeIndex)
		nodeIndex = rc.nodes[nodeIndex].prevNodeIndex
	}
	reverseUCharSlice(rc.path[rc.prefixLength:])
	rc.path = append(rc.path, 0)

	rc.value = candidate.value
	heap.Pop(&rc.candidateQueue)

	return true
}

// Pushes a node to queue.
func (rc *RankedCompleter) enqueueNode(nodeIndex baseType) {
	if rc.nodes[nodeIndex].isQueued {
		return
	}

	rc.nodeQueue = append(rc.nodeQueue, nodeIndex)
	rc.nodes[nodeIndex].isQueued = true
}

// Pushes a candidate to priority queue.
func (rc *RankedCompleter) enqueueCandidate(nodeIndex baseType) {
	heap.Push(&rc.candidateQueue, &RankedCompleterCandidate{
		nodeIndex: nodeIndex,
		value:     dictValue(rc.dict.units[rc.nodes[nodeIndex].dictIndex]),
	})
}

// Finds a sibling of a given node.
func (rc *RankedCompleter) findSibling(nodeIndex *baseType) bool {
	var prevNodeIndex baseType = rc.nodes[*nodeIndex].prevNodeIndex
	var dictIndex baseType = rc.nodes[*nodeIndex].dictIndex

	var siblingLabel ucharType = rc.guide.Sibling(dictIndex)
	if siblingLabel == 0 {
		if !rc.nodes[prevNodeIndex].hasTerminal {
			return false
		}
		rc.nodes[prevNodeIndex].hasTerminal = false
	}

	// Follows a transition to sibling and creates a node for the sibling.
	var dictPrevIndex baseType = rc.nodes[prevNodeIndex].dictIndex
	dictIndex = rc.followWithoutCheck(dictPrevIndex, siblingLabel)
	*nodeIndex = rc.createNode(dictIndex, prevNodeIndex, siblingLabel)

	return true
}

// Follows transitions and finds a terminal.
func (rc *RankedCompleter) findTerminal(nodeIndex baseType) baseType {
	for rc.nodes[nodeIndex].label != 0 {
		var dictIndex baseType = rc.nodes[nodeIndex].dictIndex
		var childLabel ucharType = rc.guide.Child(dictIndex)
		if childLabel == 0 {
			rc.nodes[nodeIndex].hasTerminal = false
		}

		// Follows a transition to child and creates a node for the child.
		dictIndex = rc.followWithoutCheck(dictIndex, childLabel)
		nodeIndex = rc.createNode(dictIndex, nodeIndex, childLabel)
	}
	return nodeIndex
}

// Follows a transition without any check.
func (rc *RankedCompleter) followWithoutCheck(index baseType, label ucharType) baseType {
	return index ^ dictOffset(rc.dict.units[index]) ^ baseType(label)
}

// Creates a node
func (rc *RankedCompleter) createNode(dictIndex baseType, prevNodeIndex baseType, label ucharType) baseType {
	rc.nodes = append(rc.nodes, RankedCompleterNode{
		dictIndex:     dictIndex,
		prevNodeIndex: prevNodeIndex,
		label:         label,
		hasTerminal:   label != 0 && rc.dict.HasValue(dictIndex),
	})
	return baseType(len(rc.nodes) - 1)
}
