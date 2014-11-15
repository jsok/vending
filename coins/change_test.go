package coins

import (
	"reflect"
	"testing"
)

func TestGreedyChanger(t *testing.T) {
	var cases = []struct {
		denoms   DenominationSlice
		value    int
		expected ChangeSet
	}{
		{
			DenominationSlice{1, 5, 10, 20, 50, 100, 200},
			45,
			ChangeSet{20: 2, 5: 1},
		},
		{
			// Typical greedy behaviour
			DenominationSlice{1, 3, 4},
			6,
			ChangeSet{4: 1, 1: 2},
		},
	}
	counter := &CoinCounter{}
	for _, c := range cases {
		greedy := NewGreedyChanger(c.denoms)
		change, err := greedy.MakeChange(c.value)
		if err != nil {
			t.Errorf("GreedyChanger could not make change for %dc", c.value)
		}
		actualValue := counter.Count(change)
		if actualValue != c.value {
			t.Errorf("GreedyChanger gave total change %dc instead of %dc", actualValue, c.value)
		}
		if !reflect.DeepEqual(change, c.expected) {
			t.Errorf("GreedyChanger made change %v, expected %v", change, c.expected)
		}

	}
}

func TestGreedyChangerFailure(t *testing.T) {
	var cases = []struct {
		denoms DenominationSlice
		value  int
	}{
		{
			DenominationSlice{3, 4},
			6,
		},
	}
	for _, c := range cases {
		greedy := NewGreedyChanger(c.denoms)
		if _, err := greedy.MakeChange(c.value); err == nil {
			t.Errorf("GreedyChanger should not be able to make exact change for %dc with denominations %v", c.value, c.denoms)
		}
	}
}
