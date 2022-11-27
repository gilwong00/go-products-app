package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"products-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, h *http.Request) {
	list := data.GetProducts()
	data, err := json.Marshal(list)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}

	w.Write(data)
}
