package machine

import (
	"fmt"
	"sort"
	"strings"
)

type Denomination int

type byDenomination []Denomination

func (v byDenomination) Len() int           { return len(v) }
func (v byDenomination) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byDenomination) Less(i, j int) bool { return v[i] < v[j] }

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
	MakeChange(value int) (Change, error)
}

type GreedyChangeMaker struct {
	denoms []Denomination
}

func NewGreedyChangeMaker(denoms []Denomination) *GreedyChangeMaker {
	// Keep the denominations sorted in descending order
	sort.Sort(sort.Reverse(byDenomination(denoms)))
	return &GreedyChangeMaker{denoms: denoms}
}

func (r *GreedyChangeMaker) MakeChange(value int) (Change, error) {
	rem := value
	change := make(Change, 0)

	for i := range r.denoms {
		d := r.denoms[i]
		for int(d) <= rem && rem > 0 {
			change[d]++
			rem -= int(d)
		}
	}

	if rem != 0 {
		return change, fmt.Errorf("Could not make exact change for value %d", value)
	}
	return change, nil
}
