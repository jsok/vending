package http

import (
	"log"
	"net/http"

	"github.com/jsok/vending/machine"
)

func Serve(machine *machine.Machine) {
	mux := http.NewServeMux()

	mux.Handle("/api/items", &itemsListHandler{machine})
	mux.Handle("/api/items/", &itemHandler{machine})
	mux.Handle("/api/purchase", &purchaseHandler{machine})

	log.Println("Listening...")
	http.ListenAndServe(":5000", mux)
}
