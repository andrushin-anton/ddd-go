// Package main runs the tavern and performs an Order
package main

import (
	"github.com/google/uuid"
	"github.com/andrushin-anton/ddd-go/domain/products"
	"github.com/andrushin-anton/ddd-go/services/order"
	servicetavern "github.com/andrushin-anton/ddd-go/services/tavern"
)

func main() {

	products := productInventory()
	// Create Order Service to use in tavern
	os, err := order.NewOrderService(
		order.WithMemoryCustomerRepository(),
		order.WithMemoryProductRepository(products),
	)
	if err != nil {
		panic(err)
	}
	// Create tavern service
	tavern, err := servicetavern.NewTavern(
		servicetavern.WithOrderService(os))
	if err != nil {
		panic(err)
	}

	uid, err := os.AddCustomer("Percy")
	if err != nil {
		panic(err)
	}
	order := []uuid.UUID{
		products[1].ID(),
	}
	// Execute Order
	err = tavern.Order(uid, order)
	if err != nil {
		panic(err)
	}
}

func productInventory() []products.Product {
	beer, err := products.NewProduct("Beer", "Healthy Beverage", 1.99)
	if err != nil {
		panic(err)
	}
	peenuts, err := products.NewProduct("Peenuts", "Healthy Snacks", 0.99)
	if err != nil {
		panic(err)
	}
	wine, err := products.NewProduct("Wine", "Healthy Snacks", 0.99)
	if err != nil {
		panic(err)
	}
	products := []products.Product{
		beer, peenuts, wine,
	}
	return products
}