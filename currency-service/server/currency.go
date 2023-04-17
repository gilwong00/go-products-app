package server

import (
	"context"

	"github.com/gilwong00/go-product/currency-service/data"
	protos "github.com/gilwong00/go-product/currency-service/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type CurrencyServer struct {
	protos.UnimplementedCurrencyServer
	rates *data.ExchangeRate
	log   hclog.Logger
}

func NewCurrencyServer(rates *data.ExchangeRate, log hclog.Logger) *CurrencyServer {
	return &CurrencyServer{rates: rates, log: log}
}

func (c *CurrencyServer) GetCurrencyRate(ctx context.Context, req *protos.GetCurrencyRateRequest) (*protos.GetCurrencyRateResponse, error) {
	c.log.Info("Handle GetCurrencyRate", "Initial", req.GetInitial(), "Output", req.GetFinal())
	rate, err := c.rates.GetRate(req.GetInitial().String(), req.GetFinal().String())
	if err != nil {
		return nil, err
	}
	return &protos.GetCurrencyRateResponse{Rate: rate}, nil
}
