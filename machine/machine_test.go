package machine

import (
	"fmt"
	"testing"

	"github.com/jsok/vending/coins"
)

type FixedItemPicker struct {
	item *Item
}

func (p *FixedItemPicker) Pick(index int) (*Item, error) {
	return p.item, nil
}

type FailingPicker struct{}

func (p *FailingPicker) Pick(index int) (*Item, error) {
	return nil, fmt.Errorf("Failed to pick item in slot %d", index)
}

var changeMaker = coins.NewGreedyChangeMaker(coins.DenominationSlice{1, 5, 10, 20, 50, 100, 200})

func TestMachineUnderpaid(t *testing.T) {
	m := NewMachine(&FixedItemPicker{&Item{"Item", 100}}, changeMaker)

	payWith := coins.Change{50: 1}
	change, err := m.Purchase(0, payWith)
	if err == nil {
		t.Errorf("Machine should have failed the purhcase")
	}

	paid := payWith.Value()
	refunded := change.Value()
	if paid > refunded {
		t.Errorf("Machine short-changed customer! Gave it %dc, only refunded %dc",
			paid, refunded)
	} else if refunded < paid {
		t.Errorf("Machine returned too much change! Gave it %dc, refunded %dc",
			paid, refunded)
	}
}

func TestMachineExactPayment(t *testing.T) {
	m := NewMachine(&FixedItemPicker{&Item{"Item", 100}}, changeMaker)

	payWith := coins.Change{50: 1, 20: 2, 10: 1}
	change, err := m.Purchase(0, payWith)
	if err != nil {
		t.Errorf("Machine should have accepted the payment")
	}

	refunded := change.Value()
	if refunded > 0 {
		t.Errorf("Machine refunded %dc although it was given the exact amount",
			refunded)
	}
}

func TestMachineExpectChange(t *testing.T) {
	m := NewMachine(&FixedItemPicker{&Item{"Item", 100}}, changeMaker)

	payWith := coins.Change{200: 1}
	change, err := m.Purchase(0, payWith)
	if err != nil {
		t.Errorf("Machine should have accepted the payment")
	}

	refunded := change.Value()
	if refunded != 100 {
		t.Errorf("Machine refunded %dc, expected change of 100c", refunded)
	}
}

func TestMachineFailedToPick(t *testing.T) {
	m := NewMachine(&FailingPicker{}, changeMaker)

	payWith := coins.Change{100: 1}
	change, err := m.Purchase(0, payWith)

	if err == nil {
		t.Errorf("Machine should have failed the purhcase")
	}

	refunded := change.Value()
	if refunded != 100 {
		t.Errorf("Machine refunded %dc, expected change of 100c", refunded)
	}
}
