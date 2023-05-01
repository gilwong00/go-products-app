package handlers

import (
	"net/http"

	"github.com/gilwong00/go-product/products-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes a product
// responses:
//
//	201: noContent

// DeleteProduct deletes a product from the database
func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := getProductID(r)
	p.l.Debug("Deleting record", "id", id)
	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Unable to delete record id does not exist")
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&data.GenericError{Message: err.Error()}, w)
		// http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		p.l.Error("Unable to delete record", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&data.GenericError{Message: err.Error()}, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
