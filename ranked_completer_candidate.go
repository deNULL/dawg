package dawg

type RankedCompleterCandidate struct {
	nodeIndex baseType
	value     valueType
}

type RankedCompleterCandidateQueue []*RankedCompleterCandidate

func makeRankedCompleterCandidateCmp(valuesCmp valueComparatorFunc) func(lhs *RankedCompleterCandidate, rhs *RankedCompleterCandidate) bool {
	return func(lhs *RankedCompleterCandidate, rhs *RankedCompleterCandidate) bool {
		if lhs.value != rhs.value {
			return valuesCmp(lhs.value, rhs.value)
		}
		return lhs.nodeIndex > rhs.nodeIndex
	}
}

func (pq RankedCompleterCandidateQueue) Len() int {
	return len(pq)
}

// TODO: this does not allow custom comparator
func (pq RankedCompleterCandidateQueue) Less(i, j int) bool {
	if pq[i].value != pq[j].value {
		return pq[i].value > pq[j].value
	}
	return pq[i].nodeIndex > pq[j].nodeIndex
}

func (pq RankedCompleterCandidateQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *RankedCompleterCandidateQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*RankedCompleterCandidate))
}

func (pq *RankedCompleterCandidateQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}
