package concolicTypes

type ConcreteValues struct {
	intVals  map[string]int
	boolVals map[string]bool
}

func newConcreteValues() *ConcreteValues {
	return &ConcreteValues{intVals: make(map[string]int), boolVals: make(map[string]bool)}
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

// ================= BOOLS =================

// initialize unseen bools to false
func (cv *ConcreteValues) getBoolValue(name string) int {
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

// // ================= UTILITY =================

func (cv *ConcreteValues) inherit(other *ConcreteValues) {
	for keyOther, valOther := range other.intVals {
		if _, seen := cv.intVals[keyOther]; !seen {
			cv.intVals[keyOther] = valOther
		}
	}

	for keyOther, valOther := range other.boolVals {
		if _, seen := cv.intVals[keyOther]; !seen {
			cv.boolVals[keyOther] = valOther
		}
	}
}
