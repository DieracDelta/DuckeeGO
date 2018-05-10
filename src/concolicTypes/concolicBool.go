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

// ================= UNOPS =================

func (self ConcolicBool) ConcBoolNot() ConcolicBool {
	return ConcolicBool{Value: !self.Value, z3Expr: self.z3Expr.Not()}
}

// ================= BINOPS =================

func ConcBoolBinopToBool(concreteFunc func(a, b bool) bool, z3Func func(az, bz z3.Bool) z3.Bool, ac, bc ConcolicBool) ConcolicBool {
	res := concreteFunc(ac.Value, bc.Value)
	sym := z3Func(ac.z3Expr, bc.z3Expr)
	return ConcolicBool{Value: res, z3Expr: sym}
}

func (self ConcolicBool) ConcBoolAnd(other ConcolicBool) ConcolicBool {
	and := func(a, b bool) bool { return a && b }
	andZ3 := func(az, bz z3.Bool) z3.Bool { return az.And(bz) }
	return ConcBoolBinopToBool(and, andZ3, self, other)
}

func (self ConcolicBool) ConcBoolOr(other ConcolicBool) ConcolicBool {
	or := func(a, b bool) bool { return a || b }
	orZ3 := func(az, bz z3.Bool) z3.Bool { return az.Or(bz) }
	return ConcBoolBinopToBool(or, orZ3, self, other)
}

func (self ConcolicBool) ConcBoolAndNot(other ConcolicBool) ConcolicBool {
	andNot := func(a, b bool) bool { return a &^ b }
	andNotZ3 := func(az, bz z3.Bool) z3.Bool { return az.And(bz.Not()) }
	return ConcBoolBinopToBool(andNot, andNotZ3, self, other)
}