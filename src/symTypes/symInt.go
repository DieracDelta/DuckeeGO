package symTypes

import "github.com/aclements/go-z3/z3"
import "reflect"
import "hash/fnv"

type SymInt struct {
	id string
}

func (self *SymInt) SymIntEquals(o interface{}) bool {
	if reflect.TypeOf(o) != SymInt {
		return false
	}
	return self.id == (SymInt(o)).id
}

func (self *SymInt) SymIntZ3Expr() z3.Int {
	return z3.Int(self.id)
}

func (self *SymInt) SymIntHash() int {
	h := fnv.New32a()
	h.Write([]byte(self.id))
	return h.Sum32()
}
