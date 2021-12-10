package memory

import (
	"testing"

	"github.com/andrushin-anton/ddd-go/domain/customers"
	"github.com/google/uuid"
)

func TestMemory_GetCustomer(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	// Create a fake customer to add to repository
	cust, err := customers.NewCustomer("Percy")
	if err != nil {
		t.Fatal(err)
	}
	id := cust.ID()
	// Create the repo to use, and add some test Data to it for testing
	// Skip Factory for this
	repo := InMemoryCustomerRepository{
		customers: map[uuid.UUID]customers.Customer{
			id: cust,
		},
	}

	testCases := []testCase{
		{
			name:        "No Customer By ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: customers.ErrCustomerNotFound,
		}, {
			name:        "Customer By ID",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemory_AddCustomer(t *testing.T) {
	type testCase struct {
		name        string
		cust        string
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "Add Customer",
			cust:        "Percy",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := InMemoryCustomerRepository{
				customers: map[uuid.UUID]customers.Customer{},
			}

			cust, err := customers.NewCustomer(tc.cust)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Add(cust)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			found, err := repo.Get(cust.ID())
			if err != nil {
				t.Fatal(err)
			}
			if found.ID() != cust.ID() {
				t.Errorf("Expected %v, got %v", cust.ID(), found.ID())
			}
		})
	}
}
