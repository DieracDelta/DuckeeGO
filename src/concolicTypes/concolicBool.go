package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
	Value 	bool
	Z3Expr 	z3.Bool
}

func MakeConcolicBoolVar(name string) ConcolicBool {
	return ConcolicBool{Value: concreteValuesGlobal.GetBoolValue(name), Z3Expr: ctx.BoolConst(name)}
}

func MakeConcolicBoolConst(value bool) ConcolicBool {
	return ConcolicBool{Value: value, Z3Expr: ctx.FromBool(value)}
}

// ================= UNOPS =================

func (self ConcolicBool) ConcBoolNot() ConcolicBool {
	return ConcolicBool{Value: !self.Value, Z3Expr: self.Z3Expr.Not()}
}

// ================= BINOPS =================

func ConcBoolBinopToBool(concreteFunc func(a, b bool) bool, z3Func func(az, bz z3.Bool) z3.Bool, ac, bc ConcolicBool) ConcolicBool {
	res := concreteFunc(ac.Value, bc.Value)
	sym := z3Func(ac.Z3Expr, bc.Z3Expr)
	return ConcolicBool{Value: res, Z3Expr: sym}
}

// TODO not equal

func (self ConcolicBool) ConcEq(other ConcolicBool) ConcolicBool {
	eq := func(a, b bool) bool { return a == b }
	eqZ3 := func(az, bz z3.Bool) z3.Bool { return az.Eq(bz) }
	return ConcBoolBinopToBool(eq, eqZ3, self, other)
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
