package internal

import (
	"context"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
)

func (c *CurrencyServer) GetCurrencyRate(ctx context.Context, req *protos.GetCurrencyRateRequest) (*protos.GetCurrencyRateResponse, error) {
	c.log.Info("Handle GetCurrencyRate", "Initial", req.GetInitial(), "Output", req.GetFinal())
	rate, err := c.rates.GetRate(req.GetInitial().String(), req.GetFinal().String())
	if err != nil {
		return nil, err
	}
	return &protos.GetCurrencyRateResponse{Rate: rate}, nil
}
