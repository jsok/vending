package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	fmt.Println(config)
}
