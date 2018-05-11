package concolicTypes

import "github.com/aclements/go-z3/z3"

type SymbolicStack struct {
	argsStack []z3.Value
	retValue z3.Value
	retValid bool
}

func newSymbolicStack() *SymbolicStack {
	// return &SymbolicStack{argsStack: make([]z3.Value, 0)}
}

func (s *SymbolicStack) pushArg(val z3.Value) {
	s.argsStack = append(s.argsStack, val)
}

func (s *SymbolicStack) popArg() z3.Value {
	v := s.argsStack[len(s.argsStack) - 1]
	s.argsStack = s.argsStack[:len(s.argsStack) - 1]
	return v
}

func (s *SymbolicStack) pushRet(val z3.Value) {
	s.retValue = val
	s.retValid = true
}

func (s *SymbolicStack) popRet() z3.Value {
	ret := s.retValue
	s.retValid = false
	return ret
}

func (s *SymbolicStack) hasRetValue() bool {
	return s.retValid
}

