package machine

import (
	"fmt"
	"sort"
	"strings"
)

type Denomination int
type DenominationSlice []int

func (d DenominationSlice) String() string {
	s := make([]string, len(d))
	for i := range d {
		s[i] = fmt.Sprintf("%dc", d[i])
	}
	return fmt.Sprintf("DenominationSlice[%v]", strings.Join(s, " "))
}

type Change map[Denomination]int

func (c Change) Value() int {
	var accum int = 0
	for denom, n := range c {
		accum += int(denom) * n
	}
	return accum
}

func (c Change) String() string {
	s := make([]string, 0)
	for denom, num := range c {
		s = append(s, fmt.Sprintf("%d x %dc", num, denom))
	}
	return fmt.Sprintf("Change[%v]", strings.Join(s, ", "))
}

type ChangeMaker interface {
	MakeChange(value int) (change Change, err error)
}

type GreedyChangeMaker struct {
	denoms DenominationSlice
}

func NewGreedyChangeMaker(denoms DenominationSlice) *GreedyChangeMaker {
	// Keep the denominations sorted in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(denoms)))
	return &GreedyChangeMaker{denoms: denoms}
}

func (r *GreedyChangeMaker) MakeChange(value int) (Change, error) {
	rem := value
	change := make(Change, 0)

	for i := range r.denoms {
		d := r.denoms[i]
		for d <= rem && rem > 0 {
			change[Denomination(d)]++
			rem -= d
		}
	}

	if rem != 0 {
		return change, fmt.Errorf("Could not make exact change for value %d", value)
	}
	return change, nil
}
