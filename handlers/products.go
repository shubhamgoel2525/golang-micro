package handlers

import (
	"log"
	"net/http"

	"github.com/shubhamgoel2525/working/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (products *Products) ServeHTTP(rw http.ResponseWriter, h *http.Request) {
	if h.Method == http.MethodGet {
		products.getProducts(rw, h)
		return
	}

	// handle an update

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (products *Products) getProducts(rw http.ResponseWriter, h *http.Request) {
	productsList := data.GetProducts()
	err := productsList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
