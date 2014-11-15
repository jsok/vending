package coins

import (
	"fmt"
	"sort"
)

type Changer interface {
	MakeChange(value int) (change ChangeSet, err error)
}

type GreedyChanger struct {
	denoms DenominationSlice
}

func NewGreedyChanger(denoms DenominationSlice) *GreedyChanger {
	// Keep the denominations sorted in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(denoms)))
	return &GreedyChanger{denoms: denoms}
}

func (r *GreedyChanger) MakeChange(value int) (ChangeSet, error) {
	rem := value
	change := make(ChangeSet, 0)

	for i := range r.denoms {
		d := r.denoms[i]
		for d <= rem && rem > 0 {
			change[Denomination(d)]++
			rem -= d
		}
	}

	if rem != 0 {
		return nil, fmt.Errorf("Could not make exact change for value %d", value)
	}
	return change, nil
}
