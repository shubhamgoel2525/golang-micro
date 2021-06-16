package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shubhamgoel2525/working/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// TODO: Clear unrequired handlers
	// helloHandler := handlers.NewHello(logger)
	// goodbyeHandler := handlers.NewGoodbye(logger)
	productsHandler := handlers.NewProducts(logger)

	serveMux := http.NewServeMux()
	// serveMux.Handle("/", helloHandler)
	// serveMux.Handle("/goodbye", goodbyeHandler)
	serveMux.Handle("/", productsHandler)

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
