package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrProductNotFound = fmt.Errorf("Product not found")

//swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"sku"`
	CreatedAt string `json:"-"` // omitting from output
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type Products []*Product

type ProductsDB struct {
	currency     protos.CurrencyClient
	log          hclog.Logger
	rates        map[string]float64
	streamClient protos.Currency_StreamCurrencyRatesClient
}

type GenericError struct {
	Message string `json:"message"`
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Coffee",
		Description: "The best drink ever",
		Price:       3.99,
		SKU:         "cofe123",
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	},
	{
		ID:          2,
		Name:        "Esspresso",
		Description: "Small and strong",
		Price:       1.99,
		SKU:         "esspr123",
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	},
}

func NewProductDB(currency protos.CurrencyClient, log hclog.Logger) *ProductsDB {
	db := &ProductsDB{currency, log, make(map[string]float64), nil}
	go db.handleCurrencyUpdates()
	return db
}

func (p *ProductsDB) handleCurrencyUpdates() {
	client, err := p.currency.StreamCurrencyRates(context.Background())
	if err != nil {
		p.log.Error("unable to subscribe to rate changes")
	}
	p.streamClient = client
	for {
		rateRequest, err := client.Recv()
		p.log.Info("received update rate from service", "final", rateRequest.GetFinal().String())
		if err != nil {
			p.log.Error("error receiving message")
		}
		p.rates[rateRequest.Final.String()] = rateRequest.Rate
	}
}

// decodes json from createProduct to match Product struct
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}
	rate, err := p.getRateForProduct(currency)
	if err != nil {
		p.log.Error("unable to get rate", "currenct", currency, "error", err)
		return nil, err
	}
	response := Products{}
	for _, product := range productList {
		p := *product
		p.Price = p.Price * float32(rate)
		response = append(response, &p)
	}
	return response, nil
}

func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	if currency == "" {
		return productList[i], nil
	}
	rate, err := p.getRateForProduct(currency)
	if err != nil {
		p.log.Error("unable to get rate", "currenct", currency, "error", err)
		return nil, err
	}
	// take a copy because productList is a reference so if we mutate the actual
	// value in the productList, we update the actual collection item instead of returning a specific update
	product := *productList[i]
	product.Price = float32(rate)
	return &product, nil
}

// AddProduct adds a new product to the database
func (p *ProductsDB) AddProduct(pr Product) {
	// get the next id in sequence
	id := productList[len(productList)-1].ID
	pr.ID = id + 1
	productList = append(productList, &pr)
}

func (p *ProductsDB) UpdateProduct(product Product) error {
	i := findIndexByID(product.ID)
	if i == -1 {
		return ErrProductNotFound
	}
	// productValues := reflect.ValueOf(product).Elem() // use .Elem since product is a reference to a struct
	// productTypes := productValues.Type()
	productList[i] = &product
	return nil
}

func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByID(id)
	if i == -1 {
		return ErrProductNotFound
	}
	//removing product from list
	productList = append(productList[:i], productList[i+1])
	return nil
}

func findIndexByID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (p *ProductsDB) getRateForProduct(finalCurrency string) (float64, error) {
	// get from cache
	if rate, ok := p.rates[finalCurrency]; ok {
		return rate, nil
	}
	request := &protos.GetCurrencyRateRequest{
		Initial: protos.Currencies(protos.Currencies_value["EUR"]),
		Final:   protos.Currencies(protos.Currencies_value[finalCurrency]),
	}
	// get initial rate
	rate, err := p.currency.GetCurrencyRate(context.Background(), request)
	if err != nil {
		// convert basic error into a rich gRPC error
		if s, ok := status.FromError(err); ok {
			metadata := s.Details()[0].(*protos.GetCurrencyRateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("unable to get rate, initial and final cannot be the same value")
			}
			return -1, fmt.Errorf("unable to get rate from currency service, inital: %s, final: %s", metadata.Initial.String(), metadata.Final.String())
		}
		return -1, err
	}
	p.rates[finalCurrency] = rate.Rate
	// sub to currency service to get updates
	p.streamClient.Send((*protos.StreamCurrencyRateRequest)(request))
	return rate.Rate, err
}

// func generateId() int {
// 	product := productList[len(productList)-1]
// 	return product.ID + 1
// }

// func findProductById(id int) (*Product, int, error) {
// 	for i, p := range productList {
// 		if p.ID == id {
// 			return p, i, nil
// 		}
// 	}
// 	return nil, -1, fmt.Errorf("Product not found")
// }
