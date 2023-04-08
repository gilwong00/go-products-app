package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
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

func GetProducts() Products {
	return productList
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

func GetProductByID(id int) (*Product, error) {
	i := findIndexByID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	return productList[i], nil
}

func findIndexByID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}
