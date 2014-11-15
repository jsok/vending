package main

type Config struct {
	Denominations []int
	Slots         []struct {
		Item struct {
			Name  string
			Price int
		}
		Inventory int
	}
}
