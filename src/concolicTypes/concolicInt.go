package concolicTypes

import "symTypes"
import "github.com/aclements/go-z3/z3"

// type ConcolicBool struct {
// 	Value bool
// 	Sym   sym.SymBool
// }

// func (concBool *ConcolicBool) equals(other ConcolicBool) ConcolicBool {
//   // strange stuff...
//   res =
//   return ConcolicBool{
//     Value : res,
//     Sym   :
//   }
// }

type ConcolicInt struct {
	Value int
	Sym   symInt.SymInt
}

func (self ConcolicInt) equals(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	switch o.(type) {
	case int:
		res := self.Value == int(o)
	case ConcolicInt:
		res := self.Value == ConcolicInt(o).Value
	default:
		return ConcolicBool(false)
	}
  return ConcolicBool{Value:res, Sym: ... }
}

func (self ConcolicInt) add(o interface{}) ConcolicInt {
  switch o.(type) {
  case int:
    res := self.Value + int(o)
  case ConcolicInt:
    res := self.Value + ConcolicInt(o).Value
  default:
    // something went very wrong.
    return nil
  }
  return ConcolicInt{Value:res, Sym:sym}
}
