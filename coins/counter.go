package coins

import (
	"fmt"
	"strings"
)

type ChangeSet map[Denomination]int

func (c ChangeSet) String() string {
	s := make([]string, 0)
	for denom, num := range c {
		s = append(s, fmt.Sprintf("%d x %dc", num, denom))
	}
	return fmt.Sprintf("ChangeSet[%v]", strings.Join(s, ", "))
}

type Counter interface {
	Count(coins ChangeSet) int
}

type CoinCounter struct{}

func (c *CoinCounter) Count(coins ChangeSet) int {
	var accum int = 0
	for denom, n := range coins {
		accum += int(denom) * n
	}
	return accum
}
