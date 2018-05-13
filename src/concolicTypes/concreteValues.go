package concolicTypes

import "fmt"

var branchNumber = 0

type ConcreteValues struct {
	intVals    map[string]int
	boolVals   map[string]bool
	stringVals map[string]string
	mapVals    map[string]map[int]int
}

func newConcreteValues() *ConcreteValues {
	return &ConcreteValues{intVals: make(map[string]int), boolVals: make(map[string]bool), stringVals: make(map[string]string), mapVals: make(map[string]map[int]int)}
}

// ================= INTS =================

// initialize unseen ints to 0
func (cv *ConcreteValues) getIntValue(name string) int {
	if _, ok := cv.intVals[name]; !ok {
		cv.intVals[name] = 0
		return 0
	}
	return cv.intVals[name]
}

func (cv *ConcreteValues) getIntMappings() map[string]int {
	return cv.intVals
}

func (cv *ConcreteValues) addIntValue(name string, value int) {
	cv.intVals[name] = value
}

// ================= STRINGS =================

// initialize unseen strings to empty
func (cv *ConcreteValues) getStringValue(name string) string {
	if _, ok := cv.stringVals[name]; !ok {
		cv.stringVals[name] = ""
		return ""
	}
	return cv.stringVals[name]
}

func (cv *ConcreteValues) getStringMappings() map[string]string {
	return cv.stringVals
}

func (cv *ConcreteValues) addStringValue(name string, value string) {
	cv.stringVals[name] = value
}

// ================= BOOLS =================

// initialize unseen bools to false
func (cv *ConcreteValues) getBoolValue(name string) bool {
	if _, ok := cv.boolVals[name]; !ok {
		cv.boolVals[name] = false
		return false
	}
	return cv.boolVals[name]
}

func (cv *ConcreteValues) getBoolMappings() map[string]bool {
	return cv.boolVals
}

func (cv *ConcreteValues) addBoolValue(name string, value bool) {
	cv.boolVals[name] = value
}

// ================= MAPS =================

// initialize unseen maps to an empty map
func (cv *ConcreteValues) getMapValue(name string) map[int]int {
	if _, ok := cv.mapVals[name]; !ok {
		cv.mapVals[name] = make(map[int]int)
		return cv.mapVals[name]
	}
	return cv.mapVals[name]
}

func (cv *ConcreteValues) getMapMappings() map[string]map[int]int {
	return cv.mapVals
}

func (cv *ConcreteValues) addMapValue(name string, value map[int]int) {
	cv.mapVals[name] = value
}

// // ================= UTILITY =================

func (cv *ConcreteValues) inherit(other *ConcreteValues) {
	for keyOther, valOther := range other.intVals {
		if _, seen := cv.intVals[keyOther]; !seen {
			cv.intVals[keyOther] = valOther
		}
	}

	for keyOther, valOther := range other.boolVals {
		if _, seen := cv.boolVals[keyOther]; !seen {
			cv.boolVals[keyOther] = valOther
		}
	}

	for keyOther, valOther := range other.stringVals {
		if _, seen := cv.stringVals[keyOther]; !seen {
			cv.stringVals[keyOther] = valOther
		}
	}
}

func (cv *ConcreteValues) printValues() {
	fmt.Printf("BRANCH %v:\r\n", branchNumber)
	branchNumber++
	fmt.Println("[")
	for key, value := range cv.intVals {
		fmt.Printf("%v -> %v,\n", key, value)
	}
	for key, value := range cv.boolVals {
		fmt.Printf("%v -> %v,\n", key, value)
		/*
		   if value {
		     fmt.Println("%v -> true,\n", key)
		   } else {
		     fmt.Println(" -> false,")
		   }
		*/
	}
	for key, value := range cv.stringVals {
		fmt.Printf("%v -> %v,\n", key, value)
	}
	for key, value := range cv.mapVals {
		str := "["
		for k1, v1 := range value {
			str += fmt.Sprintf("( %v -> %v ),", k1, v1)
		}
		str += "]"
		fmt.Printf("%v -> %v,\n", key, str)
	}
	fmt.Println("]")
	fmt.Println("")
}
