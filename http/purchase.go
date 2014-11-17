package http

import (
	"encoding/json"
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
		h.createPurchase(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *purchaseHandler) createPurchase(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	choices, ok := r.PostForm["choice"]
	if !ok {
		http.Error(w, "Required field \"choice\" missing", 400)
	}
	choice := choices[0]

	payment := UrlValuesToChange(r.PostForm["coins[]"])

	item, change, err := h.machine.Purchase(choice, payment)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	uuid, err := pseudo_uuid()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(purchaseResponse{
		Id:     uuid,
		Item:   item.Name,
		Change: change,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)
}

type purchaseResponse struct {
	Id     string
	Item   string
	Change machine.Change
}

func UrlValuesToChange(input []string) machine.Change {
	change := make(machine.Change)
	for _, coin := range input {
		if denom, err := strconv.Atoi(coin); err == nil {
			change[machine.Denomination(denom)]++
		}
	}
	return change
}
