package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
  Value bool
  Sym   SymBool
}

func (self ConcolicBool) equals(o interface{}) ConcolicBool {
  switch o.(type) {
  case bool:
    res := self.Value == bool(o)
  case ConcolicBool:
    res := self.Value == ConcolicBool(o).Value
  default:
    return ConcolicBool{Value: false, Sym: ...}
  }
  return ConcolicBool{Value: res, Sym: ...}
}


