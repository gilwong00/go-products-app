package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
	"github.com/go-playground/validator"
	"github.com/hashicorp/go-hclog"
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
	currency protos.CurrencyClient
	log      hclog.Logger
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
	return &ProductsDB{currency, log}
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSku)
	return validate.Struct(p)
}

func validateSku(fl validator.FieldLevel) bool {
	// sku looks like xxx-asdx-asdsadad
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1

}

func generateId() int {
	product := productList[len(productList)-1]
	return product.ID + 1
}

func findProductById(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, fmt.Errorf("Product not found")
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

func (p *ProductsDB) getRateForProduct(finalCurrency string) (float64, error) {
	request := &protos.GetCurrencyRateRequest{
		Initial: protos.Currencies(protos.Currencies_value["EUR"]),
		Final:   protos.Currencies(protos.Currencies_value[finalCurrency]),
	}
	rate, err := p.currency.GetCurrencyRate(context.Background(), request)
	return rate.Rate, err
}

func AddProductToList(p *Product) {
	p.ID = generateId()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProductById(id)
	if err != nil {
		return err
	}
	if pos == -1 {
		return fmt.Errorf("Product not found")
	}
	// productValues := reflect.ValueOf(product).Elem() // use .Elem since product is a reference to a struct
	// productTypes := productValues.Type()
	p.ID = id
	productList[pos] = p
	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByID(id)
	if i == -1 {
		return ErrProductNotFound
	}
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
