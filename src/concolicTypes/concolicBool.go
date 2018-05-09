package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
	Value 	bool
	z3Expr 	z3.Bool
}

func MakeConcolicBoolVar(cv *ConcreteValues, name string) ConcolicBool {
	return ConcolicBool{Value: cv.GetBoolValue(name), z3Expr: ctx.BoolConst(name)}
}

func MakeConcolicBoolConst(value bool) ConcolicBool {
	return ConcolicBool{Value: value, z3Expr: ctx.FromBool(value)}
}

func (self ConcolicBool) equals(o interface{}) ConcolicBool {
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case bool:
		res = self.Value == bool(o.(bool))
		sym = self.z3Expr.Eq(ctx.FromBool(o.(bool)))
	case ConcolicBool:
		res = self.Value == o.(ConcolicBool).Value
		sym = self.z3Expr.Eq(o.(ConcolicBool).z3Expr)
	default:
		reportError("cannot compare with == : incompatible types", self, o)
		// do something?
	}
	return ConcolicBool{Value: res, z3Expr: sym}
}
