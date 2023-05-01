package handlers

import (
	"context"
	"net/http"

	"github.com/gilwong00/go-product/products-api/data"
)

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := data.FromJSON(product, r.Body)
		if err != nil {
			p.l.Error("Deserializing product", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&data.GenericError{Message: err.Error()}, w)
			return
		}
		// validate the product
		validationErrors := p.validator.Validate(product)
		if len(validationErrors) != 0 {
			p.l.Error("Validating product", "error", validationErrors)
			// return the validation messages as an array
			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: validationErrors.Errors()}, w)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, *product)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
