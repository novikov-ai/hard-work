package main

import "fmt"

type Product struct {
	ID    string
	Name  string
	Price int
}

type Builder[T any] interface {
	WithID(string) T
	WithName(string) T
	WithPrice(int) T
	Build() *Product
}

type ProductBuilder struct {
	id    string
	name  string
	price int
}

func (b *ProductBuilder) WithID(id string) *ProductBuilder {
	b.id = id
	return b
}

func (b *ProductBuilder) WithName(name string) *ProductBuilder {
	b.name = name
	return b
}

func (b *ProductBuilder) WithPrice(price int) *ProductBuilder {
	b.price = price
	return b
}

func (b *ProductBuilder) Build() *Product {
	// Validation logic
	if b.id == "" {
		panic("ID is required")
	}
	if b.price < 0 {
		panic("Price cannot be negative")
	}

	return &Product{
		ID:    b.id,
		Name:  b.name,
		Price: b.price,
	}
}

func CreateProduct[T Builder[T]](builder T, id string) T {
	return builder.WithID(id)
}

func BuildStandardProduct[T Builder[T]](builder T) *Product {
	return builder.
		WithID("default_123").
		WithName("Standard Product").
		WithPrice(999).
		Build()
}

func main() {
	// Standard product creation
	builder := &ProductBuilder{}
	product := builder.
		WithID("prod_123").
		WithName("Premium Widget").
		WithPrice(2499).
		Build()

	fmt.Printf("Standard Product: %+v\n", *product)

	// Using generic helper
	prebuilt := BuildStandardProduct(&ProductBuilder{})
	fmt.Printf("Prebuilt Product: %+v\n", *prebuilt)

	// Using F-bounded generic function
	created := CreateProduct(&ProductBuilder{}, "gen_456").
		WithName("Generated Item").
		WithPrice(500).
		Build()
	fmt.Printf("Generated Product: %+v\n", *created)
}