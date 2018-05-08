package symTypes

import "github.com/aclements/go-z3/z3"

type SymInt struct {
	id        string
  constant  bool
}

func (self *SymInt) SymIntZ3Expr() z3.Int {
	return concolicTypes.ctx.IntConst(self.id)
}
