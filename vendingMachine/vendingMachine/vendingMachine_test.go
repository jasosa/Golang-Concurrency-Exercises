package vendingmachine_test

import (
	"github.com/go-concurrency-exercises/vendingMachine/vendingMachine"
	"testing"
)

func TestNotExistentProduct(t *testing.T) {
	vm := vendingmachine.New()
	vm.FillIt(populate())
	t.Log("Given a product id that does not exist")
	{
		t.Log("When getting the product from the vending machine")
		{
			chProd, chErr, chChange := vm.GetProduct(5, 1, 2, 5)
			select {
			case <-chErr:
				t.Log("An error should be received...OK")
			case <-chProd:
				t.Fatal("A product should never be received...WRONG")
			}
			coinsReturned := 0
			for range chChange {
				coinsReturned++
			}
			if coinsReturned == 3 {
				t.Log("And the number of coins returned should be 3...OK")
			} else {
				t.Fatalf("And the number of coins returned should be 3 but is %d...WRONG", coinsReturned)
			}
		}
	}
}

func TestExistentProduct(t *testing.T) {
	vm := vendingmachine.New()
	vm.FillIt(populate())
	prodID := 1
	t.Logf("Given the product id %d", prodID)
	{
		t.Log("When getting the product from the vending machine")
		{
			chProd, chErr, chChange := vm.GetProduct(prodID, 2, 2, 5)
			select {
			case <-chErr:
				t.Fatalf("An error should not be received...WRONG")
			case <-chProd:
				t.Logf("The product %s should be received...OK", "Chocolate Bar")
			}
			coinsReturned := 0
			for range chChange {
				coinsReturned++
			}
			if coinsReturned == 0 {
				t.Log("And the number of coins returned should be 0...OK")
			} else {
				t.Fatalf("And the number of coins returned should be 0 but is %d...WRONG", coinsReturned)
			}
		}
	}
}

func populate() map[int]vendingmachine.Item {
	prods := make(map[int]vendingmachine.Item)
	prods[1] = vendingmachine.Item{Name: "Chocolate Bar", Price: 9, Units: 10}
	prods[2] = vendingmachine.Item{Name: "Milk box", Price: 13, Units: 10}
	prods[3] = vendingmachine.Item{Name: "Candy bag", Price: 11, Units: 10}
	return prods
}
