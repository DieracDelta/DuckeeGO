package concolicint

import "symInt"
import "symBool"
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

func (self ConcolicInt) equals(o interface{}) bool {
	// return concInt.Value == other.Value
	switch o.(type) {
	case int:
		res := self.Value == int(o)
	case ConcolicInt:
		res := self.Value == ConcolicInt(o).Value
	default:
		return false
	}

}
