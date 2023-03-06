package server

import (
	"context"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type CurrencyServer struct {
	protos.UnimplementedCurrencyServer
	log hclog.Logger
}

func NewCurrency(l hclog.Logger) *CurrencyServer {
	return &CurrencyServer{log: l}
}

func (c *CurrencyServer) GetCurrencyRate(ctx context.Context, r *protos.GetCurrencyRateRequest) (*protos.GetCurrencyRateResponse, error) {
	c.log.Info("Handle GetCurrencyRate", "Initial", r.GetInitial(), "Output", r.GetFinal())
	return &protos.GetCurrencyRateResponse{Rate: 1}, nil
}
