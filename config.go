package main

type Config struct {
	Denominations []int
	Slots         map[string]struct {
		Item struct {
			Name  string
			Price int
		}
		Inventory int
	}
}
