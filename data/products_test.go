package data

import "testing"

func TestChecksValidation(t *testing.T) {
	product := &Product{Name: "Shubham", Price: 100, SKU: "asd-as-da"}

	err := product.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
