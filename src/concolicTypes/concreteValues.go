package concolicTypes

type ConcreteValues struct {
	intVals map[string]int
}

func newConcreteValues() *ConcreteValues {
	return &ConcreteValues{make(map[string]int)}
}

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

func (cv *ConcreteValues) inherit(other *ConcreteValues) {
	for keyOther, valOther := range other.intVals {
		if _, seen := cv.intVals[keyOther]; !seen {
			cv.intVals[keyOther] = valOther
		}
	}
}