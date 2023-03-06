package main

import (
	"currency-service/server"
	"net"
	"os"

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

	gs := grpc.NewServer()
	c := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, c)
	//enable reflection api
	reflection.Register(gs)

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gs.Serve(l)
}
