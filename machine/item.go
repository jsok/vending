package machine

import "fmt"

type Item struct {
	Name  string
	Price int // in cents
}

func (i *Item) String() string {
	return fmt.Sprintf("Item[%s %dc]", i.Name, i.Price)
}
