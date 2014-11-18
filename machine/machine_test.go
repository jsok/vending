package machine

import (
	"fmt"
	"testing"
)

func TestMachineFailedToMakeChange(t *testing.T) {
	m := New(NewDefaultVendor(), &AlwaysFailingChangeMaker{})
	m.Stock("A0", "Item 0", 99, 1)

	payment := Change{100: 1}
	item, change, err := m.Purchase("A0", payment)

	if err == nil {
		t.Errorf("Purchase should have failed")
	}
	if item != nil {
		t.Errorf("No item should be returned if choice out of stock")
	}
	if payment.Value() != change.Value() {
		t.Errorf("Did not receive a full refund, paid %dc, was refunded %dc",
			payment.Value(), change.Value())
	}
}

func TestMachineWithDefaultVendor(t *testing.T) {
	m := New(NewDefaultVendor(), NewAussieChangeMaker())
	m.Stock("A0", "Item 0", 99, 10)
	m.Stock("A1", "Item 1", 50, 1)
	m.Stock("A2", "Item 2", 150, 0)

	var cases = []struct {
		Choice         string
		PayWith        Change
		ShouldSucceed  bool
		ExpectedChange int
		Reason         string
	}{
		{"A0", Change{50: 1, 20: 2, 5: 1, 1: 4}, true, 0, "Exact payment"},
		{"A0", Change{50: 1}, false, 50, "Not enough payment"},

		{"A1", Change{100: 1}, true, 50, "Should succeed and receive change"},
		{"A1", Change{50: 1}, false, 50, "Previous purchase should have exhausted inventory"},

		{"A2", Change{100: 1, 50: 1}, false, 150, "Inventory is exhausted"},
		{"A3", Change{50: 1}, false, 50, "No such item"},
	}

	for _, c := range cases {
		item, change, err := m.Purchase(c.Choice, c.PayWith)
		succeeded := err == nil

		if c.ShouldSucceed && !succeeded {
			t.Errorf("Test should have succeeded because \"%s\"", c.Reason)
		} else if !c.ShouldSucceed && succeeded {
			t.Errorf("Test should have failed because \"%s\"", c.Reason)
		}

		if c.ShouldSucceed && succeeded && item == nil {
			t.Errorf("The test passed because \"%s\""+
				"but did not return an item",
				c.Reason)
		}
		if c.ExpectedChange != change.Value() {
			t.Errorf("The test passed because \"%s\""+
				" and returned %dc change"+
				" however we expected %dc change",
				c.Reason, change.Value(), c.ExpectedChange)
		}
	}
}

func TestRefill(t *testing.T) {
	m := New(NewDefaultVendor(), NewAussieChangeMaker())
	m.Stock("A0", "Item 0", 100, 0)

	if err := m.Refill("A0", 1); err != nil {
		t.Errorf("Failed to refill machine because \"%s\"", err.Error())
	}

	if _, _, err := m.Purchase("A0", Change{100: 1}); err != nil {
		t.Errorf("Failed to purchase item after refilling it")
	}
}

func TestOutOfOrder(t *testing.T) {
	m := New(NewDefaultVendor(), NewAussieChangeMaker())
	m.Stock("A0", "Item 0", 100, 1)

	if err := m.OutOfOrder("A0"); err != nil {
		t.Errorf("Failed to set slot as Out Of Order because \"%s\"", err.Error())
	}

	if _, _, err := m.Purchase("A0", Change{100: 1}); err == nil {
		t.Errorf("Purchasing an out of order item should fail")
	}
}

//////////////////////////////////////////////////////////////////////////////
// Useful stubs
//////////////////////////////////////////////////////////////////////////////

type AlwaysFailingChangeMaker struct{}

func (cm *AlwaysFailingChangeMaker) MakeChange(value int) (Change, error) {
	return nil, fmt.Errorf("Could not make change for %dc", value)
}

func NewAussieChangeMaker() ChangeMaker {
	return NewGreedyChangeMaker([]Denomination{1, 5, 10, 20, 50, 100, 200})
}
