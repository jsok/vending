package http

import (
	"encoding/json"
	"net/http"

	"github.com/jsok/vending/machine"
)

type itemsListHandler struct {
	machine *machine.Machine
}

func (h *itemsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b, err := json.Marshal(h.machine.Describe())
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write(b)
	default:
		http.NotFound(w, r)
	}
}
