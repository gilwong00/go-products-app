package internal

import (
	"io"

	protos "github.com/gilwong00/go-product/currency-service/proto/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Bidi streaming
func (c *CurrencyServer) StreamCurrencyRates(s protos.Currency_StreamCurrencyRatesServer) error {
	/*
		Recv is a blocking method. So until the client sends a message.
		This is just going to block forever. Like netconn in the standard lib
	*/
	for {
		req, err := s.Recv()
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("client closed connection")
			break
		}
		// all other errors meants the transport between server and client is unavailable
		if err != nil {
			c.log.Error("unable to read from client", "error", err)
			return err
		}
		c.log.Info("handle client request", "request", req)
		requests, ok := c.subscriptions[s]
		if !ok {
			// initialize empty collection of rate request
			requests = []*protos.StreamCurrencyRateRequest{}
		}
		// check that subscription does not exist
		for _, v := range requests {
			if v.Initial == req.Initial && v.Final == req.Final {
				// subscription already exists
				c.log.Error("subscription is already active", "initial", req.Initial.String(), "final", req.Final.String())
				grpcError := status.New(codes.InvalidArgument, "subscription is already active")
				grpcError, err = grpcError.WithDetails(req)
				if err != nil {
					c.log.Error("unable to append metadata to error message", "error", err)
					continue
				}
				// cannot return an error because that will terminate the connection
				// instead send an error message and that will be handled by the clients Recv stream.
				request := &protos.StreamCurrencyRateResponse_Error{Error: grpcError.Proto()}
				s.Send(&protos.StreamCurrencyRateResponse{Message: request})
			}
		}
		requests = append(requests, req)
		c.subscriptions[s] = requests
	}
	return nil
}
