package main

import "concolicTypes"
import "fmt"

func main() {
	x := concolicTypes.MakeFuzzyInt("x", 6)
	y := concolicTypes.MakeFuzzyBool("y", true)

	z := f(x, y)

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
