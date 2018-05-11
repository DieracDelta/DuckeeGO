package concolicTypes

import "github.com/aclements/go-z3/z3"

type SymbolicStack struct {
	argsStack 	[]z3.Value
	argsPushed 	bool
	retValue 		z3.Value
}

func newSymbolicStack() *SymbolicStack {
	return &SymbolicStack{argsStack: make([]z3.Value, 0), argsPushed: false, retValue: nil}
}

func (s *SymbolicStack) PushArg(val z3.Value) {
	s.argsStack = append(s.argsStack, val)
}

func (s *SymbolicStack) PopArg() z3.Value {
	v := s.argsStack[len(s.argsStack) - 1]
	s.argsStack = s.argsStack[:len(s.argsStack) - 1]
	return v
}

func (s *SymbolicStack) SetArgsPushed() {
	s.argsPushed = true
}

func (s *SymbolicStack) SetArgsPopped() {
	s.argsPushed = false
}

func (s *SymbolicStack) AreArgsPushed() bool {
	return s.argsPushed
}

func (s *SymbolicStack) PushReturn(val z3.Value) {
	s.retValue = val
}

func (s *SymbolicStack) PopReturn() z3.Value {
	return s.retValue
}

func (s *SymbolicStack) ClearArgs() {
	s.argsStack = make([]z3.Value, 0)
}