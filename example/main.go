package main

import "fmt"
import "concolicTypes"

// person id -> balance
// b/c strings are mean

func transfer(balances map[int]int, sender int, recipient int, zoobars int) {
	sender_balance := balances[sender] - zoobars
	recipient_balance := balances[recipient] + zoobars

	if sender_balance < 0 || recipient_balance < 0 {
		// WHAT HAPPENS ???? :O
		fmt.Println("we failedddddd")
	}

	balances[sender] = sender_balance
	balances[recipient] = recipient_balance
}

func main() {
	balances := map[int]int{}
	// balances := make(map[int]int)

	balances = concolicTypes.MakeFuzzyMapIntInt("balances", balances)

	alex := concolicTypes.MakeFuzzyInt("alex", 0)
	bobette := concolicTypes.MakeFuzzyInt("bobette", 1)

	balances[alex] = 10
	balances[bobette] = 10

	zoobars := concolicTypes.MakeFuzzyInt("zoobars", 10)

	sum := balances[alex] + balances[bobette]

	transfer(balances, alex, bobette, zoobars)

	if balances[alex]+balances[bobette] != sum {
		fmt.Println("oh no")
	}

	g()
}

func g() {
	x := concolicTypes.MakeFuzzyInt("x", 6)
	y := concolicTypes.MakeFuzzyBool("y", true)

	z := f(x, y)
	a := map[int]int{}
	a[0] = 1

	fmt.Printf("bruh %v\r\n", z)
}

func f(x int, y bool) bool {
	z := x
	if y {
		z = -z
		z = -15 * z
	}
	if z > 0 {
		fmt.Println("hi")
		return true
	} else {
		fmt.Println("I'm tired --chris")
		return false
	}
}
