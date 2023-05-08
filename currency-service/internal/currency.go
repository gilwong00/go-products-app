package internal

import (
	"time"

	"github.com/gilwong00/go-product/currency-service/data"
	protos "github.com/gilwong00/go-product/currency-service/proto/currency"

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
				// send response to client
				err = client.Send(&protos.StreamCurrencyRateResponse{
					Message: &protos.StreamCurrencyRateResponse_RateResponse{
						RateResponse: &protos.StreamRateResponse{
							Initial: rateRequest.Initial,
							Final:   rateRequest.Final,
							Rate:    rate,
						},
					},
				})
				if err != nil {
					c.log.Error("Unable to send updated rate", "initial", rateRequest.GetInitial().String(), "final", rateRequest.GetFinal().String())
				}
			}
		}
	}
}
