package concolicTypes

import "github.com/aclements/go-z3/z3"
import "math/big"

type ConcolicString struct {
	Value  string
	Z3Expr z3.BV
}

func StringToBigInt(s string) *big.Int {
	bRep := []byte(s)
	bigRep := new(big.Int)
	bigRep.SetBytes(bRep)
	return bigRep
}

func StringToZ3BV(s string) z3.BV {
	return ctx.FromBigInt(StringToBigInt(s), ctx.BVSort(StringToBVLen(s))).(z3.BV)
}

func StringToBVLen(s string) int {
	return StringToBigInt(s).BitLen()
}

func MakeConcolicStringVar(cv *ConcreteValues, name string) ConcolicString {
	value := cv.getStringValue(name)

	return ConcolicString{
		Value:  value,
		Z3Expr: ctx.BVConst(name, StringToBVLen(value))}
}

func MakeConcolicStringConst(value string) ConcolicString {
	return ConcolicString{
		Value:  value,
		Z3Expr: ctx.FromBigInt(StringToBigInt(value), ctx.BVSort(StringToBVLen(value))).(z3.BV)}
}

// ================= UNOPS =================

// // ================= BINOPS =================

// // TODO not equal

// func (self ConcolicString) ConcEq(other ConcolicString) ConcolicString {
// 	// eq := func(a, b bool) bool { return a == b }
// 	// eqZ3 := func(az, bz z3.Bool) z3.Bool { return az.Eq(bz) }
// 	// return ConcBoolBinopToBool(eq, eqZ3, self, other)
// }

// func (self ConcolicBool) ConcBoolAnd(other ConcolicBool) ConcolicBool {
// 	and := func(a, b bool) bool { return a && b }
// 	andZ3 := func(az, bz z3.Bool) z3.Bool { return az.And(bz) }
// 	return ConcBoolBinopToBool(and, andZ3, self, other)
// }

// func (self ConcolicBool) ConcBoolOr(other ConcolicBool) ConcolicBool {
// 	or := func(a, b bool) bool { return a || b }
// 	orZ3 := func(az, bz z3.Bool) z3.Bool { return az.Or(bz) }
// 	return ConcBoolBinopToBool(or, orZ3, self, other)
// }
