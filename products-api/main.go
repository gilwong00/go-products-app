package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gilwong00/go-product/products-api/data"
	"github.com/gilwong00/go-product/products-api/handlers"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	protos "github.com/gilwong00/go-product/currency-service/protos/currency"

	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	l := hclog.Default()

	//proto client - allow insecure connection for now
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	// create proto client for currency
	cc := protos.NewCurrencyClient(conn)
	// create database instance
	db := data.NewProductDB(cc, l)

	// create the handlers
	ph := handlers.NewProducts(l, cc, db)

	// create a new serve mux and register the handlers
	// standard lib approach
	// sm := http.NewServeMux()
	// sm.Handle("/", ph)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	putRouter := sm.Methods("PUT").Subrouter()
	postRouter := sm.Methods("POST").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)
	postRouter.HandleFunc("/", ph.CreateProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	// docs route
	getRouter.Handle("/docs", sh)
	// this allows the docs route to serve swagger docs
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	// for prod we have a move defined whitelist
	corsHandler := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:    ":9090",         // configure the bind address
		Handler: corsHandler(sm), // set the default handler
		// ErrorLog:     l,                 // set the logger for the server
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		// l.Println("Starting server on port 9090")
		l.Info("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			// l.Printf("Error starting server: %s\n", err)
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
