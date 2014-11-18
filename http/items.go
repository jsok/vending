package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jsok/vending/machine"
)

type itemsListHandler struct {
	*machine.Machine
}

func (h *itemsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b, err := json.Marshal(h.Describe())
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write(b)
	default:
		http.NotFound(w, r)
	}
}

type itemHandler struct {
	*machine.Machine
}

func (h *itemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		stock(w, r)
	case "PUT":
		refill(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *itemHandler) refill(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	choice := parts[len(parts)-1]

	r.ParseForm()
	inventory, ok := r.PostForm["inventory"]
	if !ok {
		http.Error(w, "Required field \"inventory\" missing", 400)
		return
	}
	amount, err := strconv.Atoi(inventory[0])
	if err != nil {
		http.Error(w, "Please specify inventory amount as an integer", 400)
		break
	}

	if err := h.Refill(choice, amount); err != nil {
		http.Error(w, fmt.Sprintf("Could not refill \"%s\" because: %v", choice, err), 400)
		return
	}

	b, err := json.Marshal(okResponse{"OK"})
	if err != nil {
		http.Error(w, "Unknown error occurred", 500)
		return
	}
	w.Write(b)
}

type okResponse struct {
	Status string `json:"status"`
}
