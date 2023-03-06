package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
	"github.com/gilwong00/go-product/products-api/data"
)

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
}

// standard lib approach
// func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.getProducts(w, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		p.createProduct(w, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		// regex to get prod id in route
// 		regex := regexp.MustCompile(`/([0-9]+)`)
// 		g := regex.FindAllStringSubmatch(r.URL.Path, -1)

// 		if len(g) != 1 {
// 			http.Error(w, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		if len(g[0]) != 2 {
// 			http.Error(w, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		idString := g[0][1]
// 		id, err := strconv.Atoi(idString)

// 		if err != nil {
// 			http.Error(w, "Invalid product id", http.StatusBadRequest)
// 			return
// 		}

// 		p.updateProduct(id, w, r)
// 		return
// 	}

// 	// catch all
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// }

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		err := product.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "unable to marshal json", http.StatusNotFound)
			return
		}

		// validate the product
		err = product.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("product failed validation: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
