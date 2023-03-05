// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

// swagger:response productResponse
type ProductResponse struct {
	// in: body
	Body struct {
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
		SKU string `json:"sku" validate:"sku"`
	}
}

// swagger:parameters updateProduct
type productIDParamsWrapper struct {
	// The id of the product for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}
