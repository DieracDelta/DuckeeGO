package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicInt struct {
	Value int
	Sym   SymInt
}

// TODO: update these with z3 contexts

func MakeConcolicIntVar(cv *ConcreteValues, name string) ConcolicInt {
	return ConcolicInt{cv.getIntValue(name), makeSymIntVar(name)}
}

func MakeConcolicIntConst(value int) {
	return ConcolicInt{Value: value, Sym: ctx.FromInt(v)}
}

// ================= BINOPS RETURNING BOOLS =================

func ConcIntBinopToBool(concreteFunc func(a, b int) bool, z3Func func(az, bz z3.Int) z3.Bool, ac, bc ConcolicInt) ConcolicBool {
	res := concreteFunc(ac.Value, bc.Value)
	sym := z3Func(ac.Sym.z3Expr, bc.Sym.z3Expr)
	return ConcolicBool{Value: res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcEq(other ConcolicInt) ConcolicBool {
	eq := func(a, b int) bool { return a == b }
	eqZ3 := func(az, bz z3.Int) z3.Bool { return az.Eq(bz) }
	return ConcIntBinopToBool(eq, eqZ3, self, other)
}

func (self ConcolicInt) ConcNE(other ConcolicInt) ConcolicBool {
	neq := func(a, b int) bool { return a != b }
	neqZ3 := func(az, bz z3.Int) z3.Bool { return az.Eq(bz).Not() }
	return ConcIntBinopToBool(neq, neqZ3, self, other)
}

func (self ConcolicInt) ConcLT(other ConcolicInt) ConcolicBool {
	lt := func(a, b int) bool { return a < b }
	ltZ3 := func(az, bz z3.Int) z3.Bool { return az.LT(bz) }
	return ConcIntBinopToBool(lt, ltZ3, self, other)
}

func (self ConcolicInt) ConcLE(other ConcolicInt) ConcolicBool {
	le := func(a, b int) bool { return a <= b }
	leZ3 := func(az, bz z3.Int) z3.Bool { return az.LE(bz) }
	return ConcIntBinopToBool(le, leZ3, self, other)
}

func (self ConcolicInt) ConcGT(other ConcolicInt) ConcolicBool {
	gt := func(a, b int) bool { return a > b }
	gtZ3 := func(az, bz z3.Int) z3.Bool { return az.GT(bz) }
	return ConcIntBinopToBool(gt, gtZ3, self, other)
}

func (self ConcolicInt) ConcGE(other ConcolicInt) ConcolicBool {
	ge := func(a, b int) bool { return a >= b }
	geZ3 := func(az, bz z3.Int) z3.Bool { return az.GE(bz) }
	return ConcIntBinopToBool(ge, geZ3, self, other)
}

// ================= BINOPS RETURNING INTS =================

func ConcIntBinopToInt(concreteFunc func(a, b int) int, z3Func func(az, bz z3.Int) z3.Int, ac, bc ConcolicInt) ConcolicInt {
	res := concreteFunc(ac.Value, bc.Value)
	sym := z3Func(ac.Sym.z3Expr, bc.Sym.z3Expr)
	return ConcolicInt{Value: res, Sym: SymInt{sym}}
}

func (self ConcolicInt) ConcAdd(other ConcolicInt) ConcolicInt {
	add := func(a, b int) int { return a + b }
	addZ3 := func(az, bz z3.Int) z3.Int { return az.Add(bz) }
	return ConcIntBinopToInt(add, addZ3, self, other)
}

func (self ConcolicInt) ConcSub(other ConcolicInt) ConcolicInt {
	sub := func(a, b int) int { return a - b }
	subZ3 := func(az, bz z3.Int) z3.Int { return az.Sub(bz) }
	return ConcIntBinopToInt(sub, subZ3, self, other)
}

func (self ConcolicInt) ConcMul(other ConcolicInt) ConcolicInt {
	mul := func(a, b int) int { return a * b }
	mulZ3 := func(az, bz z3.Int) z3.Int { return az.Mul(bz) }
	return ConcIntBinopToInt(mul, mulZ3, self, other)
}

func (self ConcolicInt) ConcDiv(other ConcolicInt) ConcolicInt {
	div := func(a, b int) int { return a / b }
	divZ3 := func(az, bz z3.Int) z3.Int { return az.Div(bz) }
	return ConcIntBinopToInt(div, divZ3, self, other)
}

func (self ConcolicInt) ConcMod(other ConcolicInt) ConcolicInt {
	mod := func(a, b int) int { return a % b }
	modZ3 := func(az, bz z3.Int) z3.Int { return az.Mod(bz) }
	return ConcIntBinopToInt(mod, modZ3, self, other)
}

/*

func (self ConcolicInt) ConcEq(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case int:
		res = self.Value == o.(int)
    sym = self.Sym.z3Expr.Eq(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value == o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.Eq(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicBool{Value:res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcNE(o interface{}) ConcolicBool {
	eqcb := self.ConcEq(o)
	return ConcolicBool{Value:!eqcb.Value, Sym: SymBool{eqcb.Sym.z3Expr.Not()}}
}

func (self ConcolicInt) ConcLT(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case int:
		res = self.Value < o.(int)
    sym = self.Sym.z3Expr.LT(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value < o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.LT(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicBool{Value:res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcLE(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case int:
		res = self.Value <= o.(int)
    sym = self.Sym.z3Expr.LE(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value <= o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.LE(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicBool{Value:res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcGT(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case int:
		res = self.Value > o.(int)
    sym = self.Sym.z3Expr.GT(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value > o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.GT(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicBool{Value:res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcGE(o interface{}) ConcolicBool {
	// return concInt.Value == other.Value
	var res bool
	var sym z3.Bool
	switch o.(type) {
	case int:
		res = self.Value >= o.(int)
    sym = self.Sym.z3Expr.GE(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value >= o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.GE(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicBool{Value:res, Sym: SymBool{sym}}
}

func (self ConcolicInt) ConcAdd(o interface{}) ConcolicInt {
	// return concInt.Value == other.Value
	var res int
	var sym z3.Int
	switch o.(type) {
	case int:
		res = self.Value + o.(int)
    sym = self.Sym.z3Expr.Add(ctx.FromInt(int64(o.(int)), ctx.IntSort()).(z3.Int))
	case ConcolicInt:
		res = self.Value + o.(ConcolicInt).Value
    sym = self.Sym.z3Expr.Add(o.(ConcolicInt).Sym.z3Expr)
	default:
    reportError("cannot compare with == : incompatible types", self, o)
    // do something?
    //return ConcolicBool{Value: false, Sym: nil}
	}
  return ConcolicInt{Value:res, Sym: SymInt{sym}}
}
*/
