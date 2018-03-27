package vendingmachine

import (
	"errors"
	"log"
	"time"
)

//VendingMachine ...
type VendingMachine struct {
	coins        chan int
	dispensor    chan string
	change       chan int
	fail         chan error
	products     map[int]Item
	allowedCoins []Coin
}

//Item represents an item in the vending machine
type Item struct {
	Name  string
	Price int
	Units int
}

//New initializes the Vending Machine
func New() *VendingMachine {
	vm := VendingMachine{
		coins:     make(chan int),
		dispensor: make(chan string),
		change:    make(chan int),
		fail:      make(chan error),
	}

	vm.allowedCoins = []Coin{
		Coin{currency: "$", value: 1},
		Coin{currency: "$", value: 2},
		Coin{currency: "$", value: 5},
	}

	return &vm
}

//FillIt fills the machine with the set of products provided.
//It replaces all current products in the machine
func (vm *VendingMachine) FillIt(products map[int]Item) {
	vm.products = products
}

//GetProduct returns:
// a channel where the product will be dispensed
// a channel to show any possible error
// a channel to return coins if the user pay more than the product price
func (vm *VendingMachine) GetProduct(productID int, coins ...int) (<-chan string, <-chan error, <-chan int) {

	prod := vm.products[productID]

	if prod == (Item{}) {
		go func() {
			vm.fail <- errors.New("Product does not exist")
			vm.returnCoins(coins) //return all the coins
		}()
	} else {

		dispenseFunc := createDispensorFunc(prod.Name, vm.dispensor)

		go vm.insertCoins(coins)

		go func() {
			vm.waitForMoneyAndDispense(prod.Price, dispenseFunc)
			vm.returnCoins(nil)
		}()
	}

	return vm.dispensor, vm.fail, vm.change
}

func (vm *VendingMachine) returnCoins(coins []int) {
	for _, coin := range coins {
		vm.change <- coin
	}
	close(vm.change)
}

/* func getCoinsToReturn(coins []int, price int){


} */

//Returns a function and a channel
//The function will dispense the product to the channel
var createDispensorFunc = func(product string, disp chan string) func() {
	return func() {
		dispense(product, disp)
	}
}

//Insert all the coins into the channel
func (vm *VendingMachine) insertCoins(coins []int) {
	for _, coin := range coins {
		log.Println("Inserting coins...", coin)
		vm.coins <- coin
		time.Sleep(500 * time.Millisecond) //to simulate some lag throwing the coins
	}

	log.Println("Closing coins channel...")
	close(vm.coins)
}

//Read coins from the channel until we get the
// amount needed for paying the item and then dispense the product
func (vm *VendingMachine) waitForMoneyAndDispense(itprice int, dispense func()) {
	amount := 0
	for {
		amount = amount + <-vm.coins
		log.Println("Receiving money...", amount)
		if amount >= itprice {
			go dispense()
			break //don't need to receive more coins
		}
	}
}

//Dispenses the product through the specified channel
func dispense(product string, ch chan string) {
	log.Println("Dispensing item...")
	ch <- product
	defer close(ch)
}
