package main

import (
	"fmt"

	entities "github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func main() {

	produto1, err := entities.NewProduct("1", "Notebook", "Computador de mesa portátil", 2500.00, 20)
	if err != nil {
		panic(err)
	}

	produto2, err := entities.NewProduct("2", "Notebook Dell", "Computador de mesa portátil", 2500.00, 20)
	if err != nil {
		panic(err)
	}

	fmt.Println(produto1.Name, produto2.Name)

}
