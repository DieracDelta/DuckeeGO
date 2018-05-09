package concolicTypes

import "github.com/aclements/go-z3/z3"

type SymInt struct {
	//id        string
  z3Expr    z3.Int
}

func makeSymIntVar(name string) SymInt {
	return SymInt{ctx.IntConst(name)}
}

func (self *SymInt) SymIntZ3Expr() z3.Int {
	return self.z3Expr//concolicTypes.ctx.IntConst(self.id)
}
