package http

import (
	"net/http"

	"github.com/jsok/vending/machine"
)

type itemsListHandler struct {
	machine *machine.Machine
}

func (h *itemsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		for _, item := range h.machine.Describe() {
			w.Write([]byte(item.String() + "\n"))
		}
	default:
		http.NotFound(w, r)
	}
}
