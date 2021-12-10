// Package memory is a in-memory implementation of the customer repository
package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/andrushin-anton/ddd-go/domain/customers"
)

// MemoryRepository fulfills the CustomerRepository interface
type InMemoryCustomerRepository struct {
	customers map[uuid.UUID]customers.Customer
	sync.Mutex
}

// NewCustomerRepository is a factory function to generate a new repository of customers
func NewCustomerRepository() *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{
		customers: make(map[uuid.UUID]customers.Customer),
	}
}

// Get finds a customer by ID
func (mr *InMemoryCustomerRepository) Get(id uuid.UUID) (customers.Customer, error) {
	if customer, ok := mr.customers[id]; ok {
		return customer, nil
	}

	return customers.Customer{}, customers.ErrCustomerNotFound
}

// Add will add a new customer to the repository
func (mr *InMemoryCustomerRepository) Add(c customers.Customer) error {
	if mr.customers == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		mr.Lock()
		mr.customers = make(map[uuid.UUID]customers.Customer)
		mr.Unlock()
	}
	// Make sure Customer isn't already in the repository
	if _, ok := mr.customers[c.ID()]; ok {
		return fmt.Errorf("customer already exists: %w", customers.ErrFailedToAddCustomer)
	}
	mr.Lock()
	mr.customers[c.ID()] = c
	mr.Unlock()
	return nil
}

// Update will replace an existing customer information with the new customer information
func (mr *InMemoryCustomerRepository) Update(c customers.Customer) error {
	// Make sure Customer is in the repository
	if _, ok := mr.customers[c.ID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", customers.ErrUpdateCustomer)
	}
	mr.Lock()
	mr.customers[c.ID()] = c
	mr.Unlock()
	return nil
}