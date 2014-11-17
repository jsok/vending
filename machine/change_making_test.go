package machine

import (
	"reflect"
	"testing"
)

func TestChangeValue(t *testing.T) {
	c := Change{1: 1, 5: 1, 10: 1, 20: 1, 50: 1}
	if c.Value() != 86 {
		t.Errorf("Value of %v should be 86", c)
	}
}

func TestGreedyChangeMaker(t *testing.T) {
	var cases = []struct {
		denoms   []Denomination
		value    int
		expected Change
	}{
		{
			[]Denomination{1},
			0,
			Change{},
		},
		{
			[]Denomination{1, 5, 10, 20, 50, 100, 200},
			45,
			Change{20: 2, 5: 1},
		},
		{
			// Typical greedy behaviour
			[]Denomination{1, 3, 4},
			6,
			Change{4: 1, 1: 2},
		},
	}
	for _, c := range cases {
		greedy := NewGreedyChangeMaker(c.denoms)
		change, err := greedy.MakeChange(c.value)
		if err != nil {
			t.Errorf("GreedyChangeMaker could not make change for %dc", c.value)
		}
		if change.Value() != c.value {
			t.Errorf("GreedyChangeMaker gave total change %dc instead of %dc", change.Value(), c.value)
		}
		if !reflect.DeepEqual(change, c.expected) {
			t.Errorf("GreedyChangeMaker made change %v, expected %v", change, c.expected)
		}

	}
}

func TestGreedyChangeMakerFailure(t *testing.T) {
	var cases = []struct {
		denoms []Denomination
		value  int
	}{
		{
			[]Denomination{3},
			2,
		},
		{
			[]Denomination{3, 4},
			6,
		},
	}
	for _, c := range cases {
		greedy := NewGreedyChangeMaker(c.denoms)
		if _, err := greedy.MakeChange(c.value); err == nil {
			t.Errorf("GreedyChangeMaker should not be able to make exact change for %dc with denominations %v", c.value, c.denoms)
		}
	}
}
