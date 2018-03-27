package main

import (
	"github.com/go-concurrency-exercises/vendingMachine/vendingMachine"
	"log"
)

func main() {

	vm := vendingmachine.New()
	vm.FillIt(populate())

	chProd, chErr, chChange := vm.GetProduct(5, 2, 2, 5)

	log.Println("Waiting for item...")

	select {
	case err := <-chErr:
		log.Println("Error received:", err.Error())
	case item := <-chProd:
		log.Println("Item dispensed:", item)
	}

	//Get the possible change
	for coin := range chChange {
		log.Println("Change returned. Coin of value:", coin)
	}
}

func populate() map[int]vendingmachine.Item {
	prods := make(map[int]vendingmachine.Item)
	prods[1] = vendingmachine.Item{Name: "Chocolate Bar", Price: 9, Units: 10}
	prods[2] = vendingmachine.Item{Name: "Milk box", Price: 13, Units: 10}
	prods[3] = vendingmachine.Item{Name: "Candy bag", Price: 11, Units: 10}
	return prods
}
