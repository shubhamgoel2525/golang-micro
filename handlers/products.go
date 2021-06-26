package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shubhamgoel2525/working/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (products *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	products.l.Println("Handle GET Products")

	productsList := data.GetProducts()

	err := productsList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (products *Products) AddProduct(rw http.ResponseWriter, h *http.Request) {
	products.l.Println("Handle POST Product")

	product := h.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
}

func (products *Products) UpdateProduct(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	products.l.Println("Handle PUT Product", id)
	product := h.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		err := product.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		request := r.WithContext(ctx)

		next.ServeHTTP(rw, request)
	})
}
