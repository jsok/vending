package machine

import (
	"fmt"
	"sort"
)

type Machine struct {
	vendor      Vendor
	changeMaker ChangeMaker
}

func New(vendor Vendor, changeMaker ChangeMaker) *Machine {
	return &Machine{vendor, changeMaker}
}

// Purchase the item in the specific slot and accept the given coins as payment
// Return success of the purchase, and associated change or a full refund in the
// event of a failure.
func (m *Machine) Purchase(choice string, payment Change) (*Item, Change, error) {
	var change Change = nil

	slot, err := m.vendor.Pick(choice)

	if err != nil {
		change = payment
		err = &ChoiceUnavailableError{choice, err.Error()}
	} else {
		paid := payment.Value()
		price := slot.Price()
		if paid < price {
			change = payment // full refund
			err = &UnderpaidError{choice, price, paid}
		} else {
			change, err = m.changeMaker.MakeChange(paid - price)
			if err != nil {
				err = &ChoiceUnavailableError{choice, err.Error()}
				change = payment
			}
		}
	}

	var item *Item = nil
	if err == nil {
		item, err = m.vendor.Dispense(slot)
		// Out of stock, out of order, etc.
		if err != nil {
			item = nil
			change = payment
		}
	}

	return item, change, err
}

func (m *Machine) Refill(choice string, amount int) error {
	slot, err := m.vendor.Pick(choice)
	if err != nil {
		err = &ChoiceUnavailableError{choice, err.Error()}
		return err
	}
	return m.vendor.Refill(slot, amount)
}

func (m *Machine) Stock(choice, name string, price, quantity int) {
	m.vendor.Stock(choice, quantity, &Item{name, price})
}

func (m *Machine) Describe() []VendingItem {
	items := make([]VendingItem, 0)
	for choice, slot := range m.vendor.List() {
		item := VendingItem{choice, slot.ItemName(), slot.Price(), slot.Available()}
		items = append(items, item)
	}
	sort.Sort(byChoice(items))
	return items
}

type VendingItem struct {
	Choice    string
	Item      string
	Price     int
	Available bool
}

func (v VendingItem) String() string {
	available := ""
	if !v.Available {
		available = " OUT OF STOCK"
	}
	return fmt.Sprintf("VendingItem[%s \"%s\" %dc%s]",
		v.Choice, v.Item, v.Price, available)
}

type byChoice []VendingItem

func (v byChoice) Len() int           { return len(v) }
func (v byChoice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byChoice) Less(i, j int) bool { return v[i].Choice < v[j].Choice }

type ChoiceUnavailableError struct {
	choice string
	reason string
}

func (e *ChoiceUnavailableError) Error() string {
	return fmt.Sprintf(
		"Sorry, your choice \"%s\" is currently unavailable for reason: %v",
		e.choice, e.reason)
}

type UnderpaidError struct {
	choice string
	price  int
	paid   int
}

func (e *UnderpaidError) Error() string {
	return fmt.Sprintf(
		"Your choice \"%s\" costs %dc, you only paid %dc. "+
			"Please insert the correct amount and try again",
		e.choice, e.price, e.paid)
}
