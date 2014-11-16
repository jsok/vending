package machine

import "fmt"

type Slot struct {
	item      *Item
	inventory int // amount of items remaining in this slot
}

func (s *Slot) String() string {
	return fmt.Sprintf("Slot[%v, %d remaining]", s.item, s.inventory)
}

func (s *Slot) Price() int {
	return s.item.Price
}

func (s *Slot) ItemName() string {
	return s.item.Name
}

func (s *Slot) Available() bool {
	return s.inventory > 0
}
