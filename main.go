package main

import "fmt"

func main() {
	moneyTransform := NewMoneyTransform(nil)
	fmt.Println(moneyTransform.ExchangeRate(VND, USD, 250000))
}
