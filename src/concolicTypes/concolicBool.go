package concolicTypes

import "symTypes"
import "github.com/aclements/go-z3/z3"

type ConcolicBool struct {
  Value bool
  Sym   symBool.
}
