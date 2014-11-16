package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jsok/vending/machine"
)

type purchaseHandler struct {
	machine *machine.Machine
}

func (h *purchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		choice := r.PostForm["choice"][0]
		change := StringToChange(r.PostForm["coins[]"])
		item, change, err := h.machine.Purchase(choice, change)
		w.Write([]byte(fmt.Sprintln("Purchase result:", item, change, err)))

	default:
		http.NotFound(w, r)
	}
}

func StringToChange(input []string) machine.Change {
	change := make(machine.Change)
	for _, coin := range input {
		denom, _ := strconv.Atoi(coin)
		change[machine.Denomination(denom)]++
	}
	return change
}
