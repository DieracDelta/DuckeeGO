package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
  Value bool
  Sym   SymBool
}

func (self ConcolicBool) equals(o interface{}) ConcolicBool {
	var res bool
	var sym z3.Bool
  switch o.(type) {
  case bool:
    res = self.Value == bool(o.(bool))
    sym = self.Sym.z3Expr.Eq(ctx.FromBool(o.(bool)))
  case ConcolicBool:
    res = self.Value == o.(ConcolicBool).Value
    sym = self.Sym.z3Expr.Eq(o.(ConcolicBool).Sym.z3Expr)
  default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
  }
  return ConcolicBool{Value: res, Sym: SymBool{sym}}
}


