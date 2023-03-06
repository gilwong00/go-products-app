package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"product-images/pkg/files"
	"product-images/pkg/handlers"
	"product-images/pkg/middleware"
	"syscall"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

const (
	port      = ":8080"
	logLevel  = "debug"
	imagePath = "./imagestore"
)

func main() {
	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString(logLevel),
		},
	)
	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	store, err := files.NewLocal(imagePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	//file handler
	fh := handlers.NewFiles(store, l)

	//middleware
	cm := middleware.Compression{}

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	postHandler := sm.Methods(http.MethodPost).Subrouter()
	getHandler := sm.Methods(http.MethodGet).Subrouter()
	postHandler.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.Upload)
	postHandler.HandleFunc("/", fh.MultiPartUpload)
	getHandler.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(imagePath))))
	getHandler.Use(cm.CompressionMiddleware)

	s := http.Server{
		Addr:         port,              // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Info("Starting image service", "port", port)
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// wait 30 seconds to finish pendingops before gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		l.Error("Server Shutdown Failed", "error", err)
	} else {
		l.Info("Server Shutdown gracefully")
	}
}
