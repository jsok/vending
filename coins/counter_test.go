package coins

import (
	"testing"
)

func TestCoinCounter(t *testing.T) {
	counter := &CoinCounter{}
	v := counter.Count(map[Denomination]int{1: 1, 5: 1, 10: 1, 20: 1, 50: 1})
	if v != 86 {
		t.Errorf("Got wrong count, Expected %d got %d", 86, v)
	}
}
