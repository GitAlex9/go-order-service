package main

import (
	"fmt"
	"log"

	"github.com/GitAlex9/go-order-service/internal/application/id"
	entities "github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func main() {

	generator := id.NewCounterGenerator()

	// productID := generator.Generate(id.ProductPrefix)

	// product, err := entities.NewProduct(
	// 	productID,
	// 	"Notebook",
	// 	"Notebook Gamer",
	// 	3500.00,
	// 	5,
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// productID1 := generator.Generate(id.ProductPrefix)

	// product1, err := entities.NewProduct(
	// 	productID1,
	// 	"Notebook",
	// 	"Notebook Gamer",
	// 	3500.00,
	// 	5,
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	customerID := generator.Generate(id.CustomerPrefix)

	customer, err := entities.NewCustomer(
		customerID,
		"Ana",
	)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	orderID := generator.Generate(id.OrderPrefix)

	order1, err := entities.NewOrderItem(orderID, "Notebook", 1500, 2)

	if err != nil {
		log.Fatal(err)
	}

	order2, err := entities.NewOrderItem(generator.Generate(id.ProductPrefix), "Celular", 1000, 2)
	if err != nil {
		log.Fatal(err)
	}

	pedidoId := generator.Generate(id.OrderPrefix)

	pedido, err := entities.NewOrder(pedidoId, customer.ID, []entities.OrderItem{*order1, *order2})

	fmt.Println(pedido.Status)
	fmt.Println(pedido.Items)
	fmt.Println(pedido.CustomerID)

	// fmt.Println("Produto criado:")
	// fmt.Printf("%+v\n\n", product)

	// fmt.Println("Produto criado:")
	// fmt.Printf("%+v\n\n", product1)

	// fmt.Println("Cliente criado:")
	// fmt.Printf("%+v\n", customer)

	// fmt.Println(generator.Generate(id.OrderPrefix))
	// fmt.Println(generator.Generate(id.OrderPrefix))
	// fmt.Println(generator.Generate(id.OrderPrefix))

	// fmt.Println(generator.Generate(id.ProductPrefix))
	// fmt.Println(generator.Generate(id.ProductPrefix))
	// fmt.Println(generator.Generate(id.ProductPrefix))

}
