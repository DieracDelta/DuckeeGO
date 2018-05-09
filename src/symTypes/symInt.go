package symTypes

import "github.com/aclements/go-z3/z3"

type SymInt struct {
	//id        string
  z3Expr    z3.Int
  constant  bool
}

func (self *SymInt) SymIntZ3Expr() z3.Int {
	return z3Expr//concolicTypes.ctx.IntConst(self.id)
}
