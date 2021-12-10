package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/andrushin-anton/ddd-go/domain/products"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]products.Product
	sync.Mutex
}

// New is a factory function to generate a new repository of customers
func NewProductsRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]products.Product),
	}
}

// GetAll returns all products as a slice
// Yes, it never returns an error, but
// A database implementation could return an error for instance
func (mpr *MemoryProductRepository) GetAll() ([]products.Product, error) {
	// Collect all Products from map
	var products []products.Product
	for _, product := range mpr.products {
		products = append(products, product)
	}
	return products, nil
}

// GetByID searches for a product based on it's ID
func (mpr *MemoryProductRepository) GetByID(id uuid.UUID) (products.Product, error) {
	if product, ok := mpr.products[uuid.UUID(id)]; ok {
		return product, nil
	}
	return products.Product{}, products.ErrProductNotFound
}

// Add will add a new product to the repository
func (mpr *MemoryProductRepository) Add(newprod products.Product) error {
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[newprod.ID()]; ok {
		return products.ErrProductAlreadyExist
	}

	mpr.products[newprod.ID()] = newprod

	return nil
}

// Update will change all values for a product based on it's ID
func (mpr *MemoryProductRepository) Update(upprod products.Product) error {
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[upprod.ID()]; !ok {
		return products.ErrProductNotFound
	}

	mpr.products[upprod.ID()] = upprod
	return nil
}

// Delete remove an product from the repository
func (mpr *MemoryProductRepository) Delete(id uuid.UUID) error {
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[id]; !ok {
		return products.ErrProductNotFound
	}
	delete(mpr.products, id)
	return nil
}