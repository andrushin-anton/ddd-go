package mongo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/andrushin-anton/ddd-go/domain/customers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCustomerRepository struct {
	db *mongo.Database
	// customer is used to store customers
	customer *mongo.Collection
}

// mongoCustomer is an internal type that is used to store a CustomerAggregate
// we make an internal struct for this to avoid coupling this mongo implementation to the customeraggregate.
// Mongo uses bson so we add tags for that
type mongoCustomer struct {
	ID   uuid.UUID `bson:"id"`
	Name string    `bson:"name"`
}

// NewFromCustomer takes in a aggregate and converts into internal structure
func NewFromCustomer(c customers.Customer) mongoCustomer {
	return mongoCustomer{
		ID:   c.ID(),
		Name: c.Name(),
	}
}

// ToAggregate converts into a customers.Customer
// this could validate all values present etc
func (m mongoCustomer) ToAggregate() customers.Customer {
	c := customers.Customer{}

	c.SetID(m.ID)
	c.SetName(m.Name)

	return c

}

// Create a new mongodb repository
func New(ctx context.Context, connectionString string) (*MongoCustomerRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// Find Metabot DB
	db := client.Database("ddd")
	customers := db.Collection("customers")

	return &MongoCustomerRepository{
		db:       db,
		customer: customers,
	}, nil
}

func (mr *MongoCustomerRepository) Get(id uuid.UUID) (customers.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := mr.customer.FindOne(ctx, bson.M{"id": id})

	var c mongoCustomer
	err := result.Decode(&c)
	if err != nil {
		return customers.Customer{}, err
	}
	// Convert to aggregate
	return c.ToAggregate(), nil
}

func (mr *MongoCustomerRepository) Add(c customers.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	internal := NewFromCustomer(c)
	_, err := mr.customer.InsertOne(ctx, internal)
	if err != nil {
		return err
	}
	return nil
}

func (mr *MongoCustomerRepository) Update(c customers.Customer) error {
	panic("to implement")
}