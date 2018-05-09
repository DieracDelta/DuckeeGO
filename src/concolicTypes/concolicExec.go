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


func concolicExecInput(testfunc reflect.Value, concreteValues *ConcreteValues) ([]reflect.Value, []z3.Bool) {
	var currPathConstrs []z3.Bool
	f := reflect.ValueOf(testfunc)
	args := []reflect.Value{ reflect.ValueOf(concreteValues), reflect.ValueOf(currPathConstrs) }
	res := f.Call(args)
	return res, currPathConstrs
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

		for b := 0; b < len(branchConstrs); b++ {
			oppConstr := concolicForceBranch(b, branchConstrs...)
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

func (h Handler) rubberducky(cv *ConcreteValues, currPathConstrs []z3.Bool) {
	var i concolicTypes.ConcolicInt
	var j concolicTypes.ConcolicInt
	i = concolicTypes.ConcolicInt{cv.getIntValue("i"), symTypes.SymInt{"i", false}}
	j = concolicTypes.ConcolicInt{cv.getIntValue("j"), symTypes.SymInt{"j", false}}
	b := i.equals(j)
	if b.value {
		currPathConstrs = append(currPathConstrs, b.Sym)
		fmt.Printf("grace is")
		b1 := i.notEquals(j)
		if b1.value {
			currPathConstrs = append(currPathConstrs, b1.Sym)
			fmt.Printf("mean")
		} else {
			currPathConstrs = append(currPathConstrs, b1.Sym.Not())
			fmt.Printf("nice")
		}
	} else {
		currPathConstrs = append(currPathConstrs, b.Sym.Not())
		fmt.Printf("ducks")
	}
	fmt.Println()
}

/*
func rubberducky() {
	var i concolicTypes.ConcolicInt

	i = concolicTypes.ConcolicInt{5, symTypes.SymInt{"i", false}}

	i = i.Add(concolicTypes.ConcolicInt{1, symTypes.SymInt{true}})

	j := concolicTypes.ConcolicInt{69, symTypes.SymInt{"j", false}}

	i = i.Sub(concolicTypes.ConcolicInt{420, symTypes.SymInt{"", true}}.Add(j))

}
*/

func main() {
	h := new(Handler)
	method := reflect.ValueOf(h).MethodByName("rubberducky")
	concolicExec(method, 100)
}
