package symTypes

import "github.com/aclements/go-z3/z3"

type SymBool struct {
  //id string
  z3Expr  z3.Bool
}

func (self *SymBool) SymBoolZ3Expr(ctx *z3.Context) z3.Bool {
  return z3Expr//ctx.BoolConst(self.id)
}

