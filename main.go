package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jsok/vending/http"
	"github.com/jsok/vending/machine"
)

func main() {
	var configFile *string = flag.String("config", "config.json", "Path to JSON config file")
	flag.Parse()

	log.Print("Starting vending machine")
	log.Printf("Attempting to load config from %s...\n", *configFile)

	config := &Config{}

	contents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(contents, config)
	if err != nil {
		log.Fatal(err)
	}

	vendor := machine.NewDefaultVendor()
	for choice, slot := range config.Slots {
		item := slot.Item
		vendor.Stock(choice, slot.Inventory, &machine.Item{item.Name, item.Price})
	}

	changeMaker := machine.NewGreedyChangeMaker(config.Denominations)

	machine := machine.NewMachine(vendor, changeMaker)

	fmt.Println("Vending Machine items available:")
	for _, item := range machine.Describe() {
		available := ""
		if !item.Available {
			available = "OUT OF STOCK"
		}
		fmt.Printf("[%s] -> %s %dc %s\n", item.Choice, item.Item, item.Price, available)
	}

	http.Serve(machine)
}
