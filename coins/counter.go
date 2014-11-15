package coins

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
