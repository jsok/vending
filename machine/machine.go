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
func (m *Machine) Purchase(slot int, payment coins.Change) (coins.Change, error) {
	item := m.picker.Pick(slot)
	paid := payment.Value()

	if paid == item.Price {
		return coins.Change{}, nil
	} else if paid < item.Price {
		return payment, fmt.Errorf("Item in slot %d costs %dc, you only paid %dc. Issuing full refund.",
			slot, item.Price, paid)
	} else {
		return m.changeMaker.MakeChange(paid - item.Price)
	}
}

type Picker interface {
	Pick(index int) *Item
}

type itemPicker struct {
	slots []Slot
}

func (p *itemPicker) Pick(index int) *Item {
	if index > len(p.slots) {
		return nil
	}
	return p.slots[index].item
}
