package concolicTypes

import "github.com/aclements/go-z3/z3"

type ConcolicMap struct {
	Value 	map[int]int
	Z3Expr 	z3.Array
}

func MakeConcolicMapVar(name string) ConcolicMap {
	return ConcolicMap{Value: concreteValuesGlobal.getMapValue(name), Z3Expr: ctx.ConstArray(ctx.ArraySort(ctx.IntSort(), ctx.IntSort()), ctx.FromInt(0, ctx.IntSort()))}
}

func MakeConcolicMapConst(value map[int]int) ConcolicMap {
  zarr := ctx.ConstArray(ctx.ArraySort(ctx.IntSort(), ctx.IntSort()), ctx.FromInt(0, ctx.IntSort()))
  for key, v := range value {
    zarr = zarr.Store(ctx.FromInt(int64(key), ctx.IntSort()), ctx.FromInt(int64(v), ctx.IntSort()))
  }

	return ConcolicMap{Value: value, Z3Expr: zarr}
}


func (self ConcolicMap) ConcMapGet(key ConcolicInt) ConcolicInt {
  res := self.Value[key.Value]
  sym := self.Z3Expr.Select(key.Z3Expr).(z3.Int)

  return ConcolicInt{Value: res, Z3Expr: sym }
}

func (self ConcolicMap) ConcMapPut(key ConcolicInt, value ConcolicInt) {
  self.Value[key.Value] = value.Value
  self.Z3Expr = self.Z3Expr.Store(key.Z3Expr, value.Z3Expr)
  // TODO: do we want to return anything guys????
}

