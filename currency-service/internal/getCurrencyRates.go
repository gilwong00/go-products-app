package internal

import (
	"context"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *CurrencyServer) GetCurrencyRate(
	ctx context.Context,
	req *protos.GetCurrencyRateRequest,
) (*protos.GetCurrencyRateResponse, error) {
	c.log.Info("Handle GetCurrencyRate", "Initial", req.GetInitial(), "Output", req.GetFinal())
	if req.Initial == req.Final {
		// sending back gRPC eeror
		// err := status.Errorf(
		// 	codes.InvalidArgument,
		// 	"initial currenct %s cannot be the same as final currency %s",
		// 	req.Initial,
		// 	req.Final,
		// )

		// creating a more detailed error
		err := status.Newf(
			codes.InvalidArgument,
			"initial currenct %s cannot be the same as final currency %s",
			req.Initial,
			req.Final,
		)
		errWithDetails, detailErr := err.WithDetails(req)
		if detailErr != nil {
			return nil, detailErr
		}
		return nil, errWithDetails.Err()
	}
	rate, err := c.rates.GetRate(req.GetInitial().String(), req.GetFinal().String())
	if err != nil {
		return nil, err
	}
	return &protos.GetCurrencyRateResponse{Rate: rate}, nil
}
