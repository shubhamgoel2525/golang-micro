package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
	if h.Method == http.MethodPost {
		products.addProduct(rw, h)
		return
	}

	// handle an item change
	if h.Method == http.MethodPut {
		// Expect an id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		group := regex.FindAllStringSubmatch(h.URL.Path, -1)

		if len(group) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Error while string conversion", http.StatusBadRequest)
			return
		}

		products.updateProduct(id, rw, h)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (products *Products) getProducts(rw http.ResponseWriter, h *http.Request) {
	products.l.Println("Handle GET Products")

	productsList := data.GetProducts()

	err := productsList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (products *Products) addProduct(rw http.ResponseWriter, h *http.Request) {
	products.l.Println("Handle POST Product")

	product := &data.Product{}

	err := product.FromJSON(h.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(product)
}

func (products *Products) updateProduct(id int, rw http.ResponseWriter, h *http.Request) {
	products.l.Println("Handle PUT Product")

	product := &data.Product{}

	err := product.FromJSON(h.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
