package tavern

import (
	"testing"

	"github.com/google/uuid"
	"github.com/andrushin-anton/ddd-go/domain/products"
	"github.com/andrushin-anton/ddd-go/services/order"
)

func init_products(t *testing.T) []products.Product {
	beer, err := products.NewProduct("Beer", "Healthy Beverage", 1.99)
	if err != nil {
		t.Error(err)
	}
	peenuts, err := products.NewProduct("Peenuts", "Healthy Snacks", 0.99)
	if err != nil {
		t.Error(err)
	}
	wine, err := products.NewProduct("Wine", "Healthy Snacks", 0.99)
	if err != nil {
		t.Error(err)
	}
	products := []products.Product{
		beer, peenuts, wine,
	}
	return products
}

func Test_Tavern(t *testing.T) {
	// Create OrderService
	products := init_products(t)

	os, err := order.NewOrderService(
		order.WithMemoryCustomerRepository(),
		order.WithMemoryProductRepository(products),
	)
	if err != nil {
		t.Error(err)
	}

	tavern, err := NewTavern(WithOrderService(os))
	if err != nil {
		t.Error(err)
	}

	uid, err := os.AddCustomer("Percy")
	if err != nil {
		t.Error(err)
	}
	
	order := []uuid.UUID{
		products[0].ID(),
	}
	// Execute Order
	err = tavern.Order(uid, order)
	if err != nil {
		t.Error(err)
	}

}

func Test_MongoTavern(t *testing.T) {
	// Create OrderService
	products := init_products(t)

	os, err := order.NewOrderService(
		order.WithMongoCustomerRepository("mongodb://localhost:27017"),
		order.WithMemoryProductRepository(products),
	)
	if err != nil {
		t.Error(err)
	}

	tavern, err := NewTavern(WithOrderService(os))
	if err != nil {
		t.Error(err)
	}

	uid, err := os.AddCustomer("Percy")
	if err != nil {
		t.Error(err)
	}
	order := []uuid.UUID{
		products[0].ID(),
	}
	// Execute Order
	err = tavern.Order(uid, order)
	if err != nil {
		t.Error(err)
	}
}
