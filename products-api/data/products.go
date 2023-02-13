package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"` // omitting from output
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type Products []*Product

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

func (p *Products) ToJSON(w io.Writer) error {
	// use Encoder instead of Marshal for slight performance benefits
	e := json.NewEncoder(w)
	return e.Encode(p)
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
