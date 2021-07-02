package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/shubhamgoel2525/working/handlers"
)

func main() {
	// Custom logger
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create a new instance of Handler (i.e. Handler on which
	// operations will be run)
	productsHandler := handlers.NewProducts(logger)

	// Manages requests
	serveMux := mux.NewRouter()

	// Subrouter registers routes sprcifically for parent
	// router
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productsHandler.GetProducts)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productsHandler.UpdateProduct)
	putRouter.Use(productsHandler.MiddlewareProductValidation)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productsHandler.AddProduct)
	postRouter.Use(productsHandler.MiddlewareProductValidation)

	// Similar http ListenAndServe, the method by default
	// uses an already created object. Here we make the server
	// from scratch
	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received Terminate, graceful shutdown", sig)

	server.ListenAndServe()

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
