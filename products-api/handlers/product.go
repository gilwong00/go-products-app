package handlers

import (
	"net/http"
	"strconv"

	"github.com/gilwong00/go-product/products-api/data"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type KeyProduct struct{}

type Products struct {
	// l              *log.Logger
	l         hclog.Logger
	validator *data.Validation
	productDB *data.ProductsDB
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewProductsHandler(
	log hclog.Logger,
	validator *data.Validation,
	productDB *data.ProductsDB,
) *Products {
	return &Products{log, validator, productDB}
}

func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}
	return id
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
