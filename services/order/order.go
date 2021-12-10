// Package services holds all the services that connects repositories into a business flow
package order

import (
	"context"
	"log"

	"github.com/andrushin-anton/ddd-go/domain/customers"
	"github.com/andrushin-anton/ddd-go/domain/products"
	"github.com/andrushin-anton/ddd-go/infrastructure/memory"
	"github.com/andrushin-anton/ddd-go/infrastructure/mongo"
	"github.com/google/uuid"
)

// OrderConfiguration is an alias for a function that will take in a pointer to an OrderService and modify it
type OrderConfiguration func(os *OrderService) error

// OrderService is a implementation of the OrderService
type OrderService struct {
	customers customers.CustomerRepository
	products  products.ProductRepository
}

// NewOrderService takes a variable amount of OrderConfiguration functions and returns a new OrderService
// Each OrderConfiguration will be called in the order they are passed in
func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	// Create the orderservice
	os := &OrderService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// CreateOrder will chain together all repositories to create an order for a customer
// will return the collected price of all Products
func (o *OrderService) CreateOrder(customerID uuid.UUID, productIDs []uuid.UUID) (float64, error) {
	// Get the customer
	c, err := o.customers.Get(customerID)
	if err != nil {
		return 0, err
	}

	// Get each Product
	var products []products.Product
	var price float64
	for _, id := range productIDs {
		p, err := o.products.GetByID(id)
		if err != nil {
			return 0, err
		}
		products = append(products, p)
		price += p.Price()
	}

	// All Products exists in store, now we can create the order
	log.Printf("Customer: %s has ordered %d products", c.Name(), len(products))

	return price, nil
}

// AddCustomer will add a new customer and return the customerID
func (o *OrderService) AddCustomer(name string) (uuid.UUID, error) {
	c, err := customers.NewCustomer(name)
	if err != nil {
		return uuid.Nil, err
	}
	// Add to Repo
	err = o.customers.Add(c)
	if err != nil {
		return uuid.Nil, err
	}

	return c.ID(), nil
}

// WithCustomerRepository applies a given customer repository to the OrderService
func WithCustomerRepository(cr customers.CustomerRepository) OrderConfiguration {
	// return a function that matches the OrderConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

// WithMemoryCustomerRepository applies a memory customer repository to the OrderService
func WithMemoryCustomerRepository() OrderConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	cr := memory.NewCustomerRepository()
	return WithCustomerRepository(cr)
}

// WithMemoryProductRepository adds an in memory product repo and adds all input products
func WithMemoryProductRepository(products []products.Product) OrderConfiguration {
	return func(os *OrderService) error {
		// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
		pr := memory.NewProductsRepository()

		// Add Items to repo
		for _, p := range products {
			err := pr.Add(p)
			if err != nil {
				return err
			}
		}
		os.products = pr
		return nil
	}
}

func WithMongoCustomerRepository(connectionString string) OrderConfiguration {
	return func(os *OrderService) error {
		// Create the mongo repo, if we needed parameters, such as connection strings they could be inputted here
		cr, err := mongo.New(context.Background(), connectionString)
		if err != nil {
			return err
		}
		os.customers = cr
		return nil
	}
}
