package internal

import (
	"io"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
)

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
