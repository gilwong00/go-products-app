package handlers

import (
	"net/http"

	"github.com/gilwong00/go-product/products-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//
//	200: productsResponse
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Debug("[GetProducts handler]")
	w.Header().Add("Content-Type", "application/json")
	currency := r.URL.Query().Get("currency")
	products, err := p.productDB.GetProducts(currency)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&data.GenericError{Message: err.Error()}, w)
		return
	}
	err = data.ToJSON(products, w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// GetProduct handles GET requests
func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	currency := r.URL.Query().Get("currency")
	id := getProductID(r)
	p.l.Debug("[DEBUG] get record id", id)
	prod, err := p.productDB.GetProductByID(id, currency)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Debug("[ERROR] fetching product", err)
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&data.GenericError{Message: err.Error()}, w)
		return
	default:
		p.l.Debug("[ERROR] fetching product", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&data.GenericError{Message: err.Error()}, w)
		return
	}
	err = data.ToJSON(prod, w)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Debug("[ERROR] serializing product", err)
	}
}
