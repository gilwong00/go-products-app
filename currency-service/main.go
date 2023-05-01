package main

import (
	"net"
	"os"

	"github.com/gilwong00/go-product/currency-service/data"
	internal "github.com/gilwong00/go-product/currency-service/internal"
	protos "github.com/gilwong00/go-product/currency-service/protos/currency"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":5000"
)

func main() {
	log := hclog.Default()
	grpcService := grpc.NewServer()
	rates, err := data.NewExchangeRange(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}
	c := internal.NewCurrencyServer(rates, log)
	protos.RegisterCurrencyServer(grpcService, c)
	// enable reflection api
	// reflection allows us to list all the rpc methods our currency service has
	// grpcurl --plaintext localhost:5000 list
	reflection.Register(grpcService)
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	err = grpcService.Serve(l)
	if err != nil {
		log.Error("unable to start grpc service", "error", err)
		os.Exit(1)
	}
}
