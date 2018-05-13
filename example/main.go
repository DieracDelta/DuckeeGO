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
    panic("we failedddddd")
  }

  balances[sender] = sender_balance
  balances[recipient] = recipient_balance
}

func main() {
  balances := map[int]int{}
  // balances := make(map[int]int)

  alex := concolicTypes.MakeFuzzyInt("alex", 0)
  bobette := concolicTypes.MakeFuzzyInt("bobette", 1)

  balances[alex] = 10
  balances[bobette] = 10

  balances = concolicTypes.MakeFuzzyMapIntInt("balances", balances)
  zoobars := concolicTypes.MakeFuzzyInt("zoobars", 10)

  sum := balances[alex] + balances[bobette]

  transfer(alex, bobette, zoobars)
  // transfer(alex, bobette, zoobars)

  if balances[alex] + balances[bobette] != sum {
    fmt.Println("oh no")
  }
}

/*
package main

import "concolicTypes"
import "fmt"

func main() {
	x := concolicTypes.MakeFuzzyInt("x", 6)
	y := concolicTypes.MakeFuzzyBool("y", true)

	z := f(x, y)
	a := map[int]int{}
	a[0] = 1
	// h := func() int {
	// 	// symStack.PushArg(j.Z3Expr)
	// 	// symStack.PushArg(i.Z3Expr)
	// 	// symStack.SetArgsPushed()
	// 	return rubberducky(i, j)
	// }()

	fmt.Printf("bruh %v\r\n", z)
}

func f(x int, y int) bool {
	z := x + y
	if z > 0 {
		fmt.Println("hi")
		return true
	} else {
		fmt.Println("I'm tired --chris")
		return false
	}

}
*/