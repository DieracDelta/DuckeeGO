package concolicTypes

import "github.com/aclements/go-z3/z3"

type SymBool struct {
  id string
}

func (self *SymBool) SymBoolZ3Expr(ctx *z3.Context) z3.Bool {
  return ctx.BoolConst(self.id)
}

