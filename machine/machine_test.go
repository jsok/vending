package machine

import (
	"fmt"
	"testing"

	"github.com/jsok/vending/coins"
)

type FixedItemPicker struct {
	item *Item
}

func (p *FixedItemPicker) Pick(id string) (*Item, error) {
	return p.item, nil
}

type FailingPicker struct{}

func (p *FailingPicker) Pick(id string) (*Item, error) {
	return nil, fmt.Errorf("Failed to pick item in slot %d", id)
}

var changeMaker = coins.NewGreedyChangeMaker(coins.DenominationSlice{1, 5, 10, 20, 50, 100, 200})

func TestMachineUnderpaid(t *testing.T) {
	m := NewMachine(&FixedItemPicker{&Item{"Item", 100}}, changeMaker)

	payWith := coins.Change{50: 1}
	change, err := m.Purchase("A0", payWith)
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
	change, err := m.Purchase("A0", payWith)
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
	change, err := m.Purchase("A0", payWith)
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
	change, err := m.Purchase("A0", payWith)

	if err == nil {
		t.Errorf("Machine should have failed the purhcase")
	}

	refunded := change.Value()
	if refunded != 100 {
		t.Errorf("Machine refunded %dc, expected change of 100c", refunded)
	}
}

func TestMachineWithItemPicket(t *testing.T) {
	picker := &ItemPicker{map[string]*Slot{
		"A0": &Slot{&Item{"Item 0", 99}, 10},
		"A1": &Slot{&Item{"Item 1", 50}, 1},
		"A2": &Slot{&Item{"Item 2", 150}, 0},
	}}
	m := NewMachine(picker, changeMaker)

	var cases = []struct {
		Slot    string
		PayWith coins.Change
		Success bool
		Reason  string
	}{
		{"A0", coins.Change{100: 1}, true, "Should succeed with change"},
		{"A1", coins.Change{50: 1}, true, "Exact payment"},

		{"A0", coins.Change{50: 1}, false, "Not enough payment"},
		{"A2", coins.Change{50: 1}, false, "Inventory is exhausted"},
		{"A1", coins.Change{50: 1}, false, "Previous purchase should have exhausted inventory"},
		{"A3", coins.Change{50: 1}, false, "No such item"},
		{"Z9", coins.Change{50: 1}, false, "No such item"},
	}

	for _, c := range cases {
		_, err := m.Purchase(c.Slot, c.PayWith)

		if c.Success && err != nil {
			t.Errorf("Machine test should have succeeded because %s", c.Reason)
		} else if !c.Success && err == nil {
			t.Errorf("Machine test should have failed because %s", c.Reason)
		}
	}
}
