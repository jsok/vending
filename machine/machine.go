package machine

import (
	"fmt"

	"github.com/jsok/vending/coins"
)

type Machine struct {
	picker      Picker
	changeMaker coins.ChangeMaker
}

func NewMachine(picker Picker, changeMaker coins.ChangeMaker) *Machine {
	return &Machine{picker, changeMaker}
}

// Purchase the item in the specific slot and accept the given coins as payment
// Return success of the purchase, and associated change or a full refund in the
// event of a failure.
func (m *Machine) Purchase(slot string, payment coins.Change) (coins.Change, error) {
	item, err := m.picker.Pick(slot)

	if err != nil {
		return payment, fmt.Errorf("Failed for reason: %v. Issuing full refund", err)
	}

	paid := payment.Value()
	if paid == item.Price {
		return coins.Change{}, nil
	} else if paid < item.Price {
		return payment, fmt.Errorf("Item in slot %d costs %dc, you only paid %dc. Issuing full refund.",
			slot, item.Price, paid)
	}

	return m.changeMaker.MakeChange(paid - item.Price)
}

// Interface which allows items to be picked from the machine for the purpose of
// selling items.
type Picker interface {
	Pick(id string) (*Item, error)
}

type ItemPicker struct {
	slots map[string]*Slot
}

func (p *ItemPicker) Pick(id string) (*Item, error) {
	slot, ok := p.slots[id]
	if !ok {
		return nil, fmt.Errorf("There are no items in slot %d", id)
	}

	if slot.inventory <= 0 {
		return nil, fmt.Errorf("The item in slot %d is out of stock", id)
	}
	slot.inventory -= 1
	return slot.item, nil
}
