package handlers

import (
	"net/http"
	"products-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[CreateProduct handler]")

	product := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Product created: %#v", product)
	data.AddProductToList(&product)
}
