package concolicTypes

import "fmt"
import "reflect"
import "github.com/aclements/go-z3/z3"


var ctx *z3.Context

func setGlobalContext() {
  ctxConfig := z3.NewContextConfig()
  ctxConfig.SetUint("timeout", 5000)
  ctx = z3.NewContext(ctxConfig)
}


func concolicExecInput(testfunc reflect.Value, concreteValues *ConcreteValues) ([]reflect.Value, *[]z3.Bool) {
	var currPathConstrs []z3.Bool
	// f := reflect.ValueOf(testfunc)
	args := []reflect.Value{ reflect.ValueOf(concreteValues), reflect.ValueOf(&currPathConstrs) }
	res := testfunc.Call(args)
	return res, &currPathConstrs
}

func concolicForceBranch(branchNum int, branchConds ...z3.Bool) z3.Bool {
	if (branchNum < len(branchConds)) {
		cond := ctx.FromBool(true).And(branchConds[0:branchNum]...).And(branchConds[branchNum].Not())
		return cond
	} else {
		return ctx.FromBool(true)
	}
}

func concolicFindInput(constraint z3.Bool, names *ConcreteValues) (bool, *ConcreteValues) {
	solver := z3.NewSolver(ctx)
	solver.Assert(constraint)
	sat, err := solver.Check()
	newInput := newConcreteValues()
	if sat {
		model := solver.Model()
		for key, _ := range (names.getIntMappings()) {
			modelValue := model.Eval(ctx.IntConst(key), true)
			if modelValue != nil {
				value, isLiteral, ok := modelValue.(z3.Int).AsInt64()
				if isLiteral && ok {
					newInput.addIntValue(key, int(value))
				}
			}
		}
		return true, newInput
	} else if err != nil {
		panic(err)
	}
	return false, newInput
}

func concolicExec(testfunc reflect.Value, maxiter int) {
	// seenAlready := make(map[*z3.Bool]bool)
	inputs := initialConcreteValueQueue()
	iter := 0
  setGlobalContext()
	// ctxConfig := z3.NewContextConfig()
	// ctxConfig.SetUint("timeout", 5000)
	// ctx := z3.NewContext(ctxConfig)
	for (iter < maxiter) && !(inputs.isEmpty()) {
		iter += 1
		inputThisTime := inputs.dequeue()
		_, branchConstrs := concolicExecInput(testfunc, inputThisTime)
		// fmt.Printf(branchConstrs.AsAST().String())
		for b := 0; b < len(*branchConstrs); b++ {
			oppConstr := concolicForceBranch(b, *branchConstrs...)
			// fmt.Printf(oppConstr.AsAST().String())
			// if _, seen := seenAlready[oppConstr]; !seen {
				// seenAlready[oppConstr] = true
				newInputFound, newInput := concolicFindInput(oppConstr, inputThisTime)
				if newInputFound {
					newInput.inherit(inputThisTime)
					inputs.enqueue(newInput)
				}
			// }
		}
	}
}

type Handler struct {}

func (h Handler) Rubberducky(cv *ConcreteValues, currPathConstrs *[]z3.Bool) int {
	var i ConcolicInt
	var j ConcolicInt
	i = makeConcolicIntVar(cv, "i")
	j = makeConcolicIntVar(cv, "j")
	b := i.equals(j)
	if b.Value {
		*currPathConstrs = append(*currPathConstrs, b.Sym.z3Expr)
		fmt.Printf("grace is ")
		b1 := i.notEquals(j)
		if b1.Value {
			*currPathConstrs = append(*currPathConstrs, b1.Sym.z3Expr)
			fmt.Printf("mean")
		} else {
			*currPathConstrs = append(*currPathConstrs, b1.Sym.z3Expr.Not())
			fmt.Printf("very helpful")
		}
	} else {
		*currPathConstrs = append(*currPathConstrs, b.Sym.z3Expr.Not())
		fmt.Printf("ducks")
	}
	fmt.Println()
	return 0
}

func Main() {
	h := new(Handler)
	method := reflect.ValueOf(h).MethodByName("Rubberducky")
	concolicExec(method, 100)
}
