package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
	Value bool
	Sym   SymBool
}

func MakeConcolicBoolVar(cv *ConcreteValues, name string) ConcolicBool {
	return ConcolicBool{Value: cv.GetBoolValue(name), Sym: makeSymBoolVar(name)}
}

func MakeConcolicBoolConst(value bool) {
	return ConcoliBool{Value: value, Sym: ctx.FromBool(name)}
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
