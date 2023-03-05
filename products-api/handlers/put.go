package handlers

import (
	"net/http"
	"products-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[UpdateProduct handler]")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

}
