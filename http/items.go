package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		h.stock(w, r)
	case "PUT":
		h.refill(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *itemHandler) stock(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	choice := parts[len(parts)-1]

	r.ParseForm()

	name, ok := r.PostForm["name"]
	if !ok {
		http.Error(w, fmt.Sprintf("Required field \"name\" is missing"), 400)
		return
	}
	price, err := intPostForm(r.PostForm, "price")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	inventory, err := intPostForm(r.PostForm, "inventory")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	h.Stock(choice, name[0], price, inventory)

	b, err := json.Marshal(okResponse{"OK"})
	if err != nil {
		http.Error(w, "Unknown error occurred", 500)
		return
	}
	w.Write(b)
}

func (h *itemHandler) refill(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	choice := parts[len(parts)-1]

	r.ParseForm()
	inventory, err := intPostForm(r.PostForm, "inventory")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := h.Refill(choice, inventory); err != nil {
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

func intPostForm(postForm url.Values, name string) (int, error) {
	fields, ok := postForm[name]
	if !ok {
		return 0, fmt.Errorf("Required field \"%s\" is missing", name)
	}
	field, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("Please specify \"%s\" as an integer", name)
	}
	return field, nil
}
