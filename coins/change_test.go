package coins

import (
	"fmt"
	"testing"
)

func TestGreedyChange(t *testing.T) {
	denoms := DenominationSlice{1, 5, 10, 20, 50, 100, 200}
	c := NewGreedyChanger(denoms)
	change, err := c.MakeChange(45)
	fmt.Println(denoms)
	fmt.Println(change)
	if err != nil {
		t.Errorf("GreedyChanger could not make change for 45c")
	}
}
