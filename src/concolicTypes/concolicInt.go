package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicInt struct {
	Value     int
	Sym       SymInt
}

// TODO: update these with z3 contexts

func makeConcolicIntVar(cv *ConcreteValues, name string) ConcolicInt {
	return ConcolicInt{cv.getIntValue(name), makeSymIntVar(name)}
}

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




