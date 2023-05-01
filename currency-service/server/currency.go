package server

import (
	"context"
	"io"
	"time"

	"github.com/gilwong00/go-product/currency-service/data"
	protos "github.com/gilwong00/go-product/currency-service/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type CurrencyServer struct {
	protos.UnimplementedCurrencyServer
	rates         *data.ExchangeRate
	log           hclog.Logger
	subscriptions map[protos.Currency_StreamCurrencyRatesServer][]*protos.StreamCurrencyRateRequest
}

func NewCurrencyServer(rates *data.ExchangeRate, log hclog.Logger) *CurrencyServer {
	server := &CurrencyServer{rates: rates, log: log, subscriptions: make(map[protos.Currency_StreamCurrencyRatesServer][]*protos.StreamCurrencyRateRequest)}
	go server.handleCurrencyUpdates()
	return server
}

func (c *CurrencyServer) handleCurrencyUpdates() {
	updated := c.rates.PollRates(5 * time.Second)
	for range updated {
		c.log.Info("got updated rates")
		// loop over subscribed clients
		for client, value := range c.subscriptions {
			// loop over rate requests
			for _, rateRequest := range value {
				rate, err := c.rates.GetRate(rateRequest.Initial.String(), rateRequest.Final.String())
				if err != nil {
					c.log.Error("unable to get updated rate", err)
				}
				client.Send(&protos.StreamCurrencyRateResponse{
					Initial: rateRequest.Initial,
					Final:   rateRequest.Final,
					Rate:    rate,
				})
			}
		}
	}

}

func (c *CurrencyServer) GetCurrencyRate(ctx context.Context, req *protos.GetCurrencyRateRequest) (*protos.GetCurrencyRateResponse, error) {
	c.log.Info("Handle GetCurrencyRate", "Initial", req.GetInitial(), "Output", req.GetFinal())
	rate, err := c.rates.GetRate(req.GetInitial().String(), req.GetFinal().String())
	if err != nil {
		return nil, err
	}
	return &protos.GetCurrencyRateResponse{Rate: rate}, nil
}

// Bidi streaming
func (c *CurrencyServer) StreamCurrencyRates(s protos.Currency_StreamCurrencyRatesServer) error {
	/*
		Recv is a blocking method. So until the client sends a message.
		This is just going to block forever. Like netconn in the standard lib
	*/
	for {
		req, err := s.Recv()
		if err == io.EOF {
			c.log.Info("client closed connection")
			break
		}
		if err != nil {
			c.log.Error("unable to read from client", "error", err)
			return err
		}
		c.log.Info("handle client request", "request", req)
		request, ok := c.subscriptions[s]
		if !ok {
			// initialize empty collection of rate request
			request = []*protos.StreamCurrencyRateRequest{}
		}
		request = append(request, req)
		c.subscriptions[s] = request
	}
	return nil
}
