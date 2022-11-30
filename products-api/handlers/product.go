package handlers

import (
	"log"
	"net/http"
	"products-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
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

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[GetProducts handler]")
	list := data.GetProducts()
	err := list.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[CreateProduct handler]")

	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Product created: %#v", product)
	data.AddProductToList(product)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[UpdateProduct handler]")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	product := &data.Product{}
	err = product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to marshal json", http.StatusNotFound)
		return
	}

	err = data.UpdateProduct(id, product)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

}
