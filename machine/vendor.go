package machine

import (
	"fmt"
)

type Vendor interface {
	List() map[string]Slot
	Pick(choice string) (*Slot, error)
	Dispense(slot *Slot) (*Item, error)
	Stock(slot string, quantity int, item *Item)
	Refill(slot *Slot, quantity int) error
	Remove(choice string) error
}

type DefaultVendor struct {
	slots map[string]*Slot
}

func NewDefaultVendor() *DefaultVendor {
	return &DefaultVendor{make(map[string]*Slot)}
}

func (v *DefaultVendor) List() map[string]Slot {
	list := make(map[string]Slot, len(v.slots))
	for choice, slot := range v.slots {
		list[choice] = *slot
	}
	return list
}

func (v *DefaultVendor) Pick(choice string) (*Slot, error) {
	slot, ok := v.slots[choice]
	if !ok {
		return nil, fmt.Errorf("No item at choice \"%s\"", choice)
	}
	return slot, nil
}

func (v *DefaultVendor) Dispense(slot *Slot) (*Item, error) {
	if !slot.Available() {
		return nil, fmt.Errorf("The item \"%s\" is out of stock", slot.item.Name)
	}
	slot.inventory -= 1
	return slot.item, nil
}

func (v *DefaultVendor) Stock(slot string, quantity int, item *Item) {
	v.slots[slot] = &Slot{item: item, inventory: quantity}
}

func (v *DefaultVendor) Refill(slot *Slot, quantity int) error {
	slot.inventory += quantity
	return nil
}

func (v *DefaultVendor) Remove(slot string) error {
	if _, err := v.Pick(slot); err != nil {
		return err
	}
	delete(v.slots, slot)
	return nil
}
