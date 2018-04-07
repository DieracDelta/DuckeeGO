package concolic

import sym
import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
  Value   bool
  Sym     sym.SymBool
}

func (concBool *ConcolicBool) equals(other ConcolicBool) ConcolicBool {
  // strange stuff...
  res = 
  return ConcolicBool{
    Value : res,
    Sym   : 
  }
}


type ConcolicInt struct {
  Value   int
  Sym     sym.SymInt
}


func (concInt *ConcolicInt) equals(other ConcolicInt) ConcolicBool {
  return concInt.Value == other.Value
}




