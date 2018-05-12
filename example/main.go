package main

import "concolicTypes"
import "fmt"

func main() {
	x := concolicTypes.MakeFuzzyInt("x", 6)
	y := 7

	z := f(x, y)

	// h := func() int {
	// 	// symStack.PushArg(j.Z3Expr)
	// 	// symStack.PushArg(i.Z3Expr)
	// 	// symStack.SetArgsPushed()
	// 	return rubberducky(i, j)
	// }()

	fmt.Printf("bruh %v\r\n", z)
}

func f(x int, y int) int {
	z := x + y
	if z > 0 {
		print("hi")
		return 1
	} else {
		print("I'm tired --chris")
		return z
	}

}
