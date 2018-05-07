package symTypes

import "github.com/aclements/go-z3/z3"
import "reflect"
import "hash/fnv"

type SymBool struct {
  id bool
}

// TODO return concolic bool?
func (self *SymBool) SymBoolEquals(o interface{}) bool {
  if reflect.TypeOf(o) != SymBool {
    return false
  }
  return self.id == (SymBool(o)).id
}

func (self *SymBool) SymBoolZ3Expr() z3.Bool {
  return z3.Bool(self.id)
}

func (self *SymBool) SymBoolHash() int {
  h := fnv.New32a()
  h.Write([]byte(self.id))
  return h.Sum32()
}

