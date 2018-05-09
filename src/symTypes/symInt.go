package symTypes

import "github.com/aclements/go-z3/z3"

type SymInt struct {
	id string
}

func (self *SymInt) SymIntZ3Expr(ctx *z3.Context) z3.Int {
	return ctx.IntConst(self.id)
}
