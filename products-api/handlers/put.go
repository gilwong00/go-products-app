package handlers

import (
	"net/http"
	"strconv"

	"github.com/gilwong00/go-product/products-api/data"
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
	w.Header().Add("Content-Type", "application/json")
	p.l.Debug("[UpdateProduct handler]")
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}
	product := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Updating record id", product.ID)
	err = p.productDB.UpdateProduct(product)
	if err == data.ErrProductNotFound {
		p.l.Error("Product not found", err)
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&data.GenericError{Message: "Product not found in database"}, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
