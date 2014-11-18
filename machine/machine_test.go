package machine

import (
	"fmt"
	"testing"
)

func TestMachineFailedToMakeChange(t *testing.T) {
	vendor := NewDefaultVendor()
	vendor.Stock("A0", 1, &Item{"Item 0", 99})
	m := New(vendor, &AlwaysFailingChangeMaker{})

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
	vendor := NewDefaultVendor()
	vendor.Stock("A0", 10, &Item{"Item 0", 99})
	vendor.Stock("A1", 1, &Item{"Item 1", 50})
	vendor.Stock("A2", 0, &Item{"Item 2", 150})

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

	m := New(vendor, NewAussieChangeMaker())

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
	vendor := NewDefaultVendor()
	vendor.Stock("A0", 0, &Item{"Item 0", 100})

	m := New(vendor, NewAussieChangeMaker())

	if err := m.Refill("A0", 1); err != nil {
		t.Errorf("Failed to refill machine because \"%s\"", err.Error())
	}

	if _, _, err := m.Purchase("A0", Change{100: 1}); err != nil {
		t.Errorf("Failed to purchase item after refilling it")
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
