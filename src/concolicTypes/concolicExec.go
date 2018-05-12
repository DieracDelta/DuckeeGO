package concolicTypes

import "fmt"
import "reflect"
import "github.com/aclements/go-z3/z3"
import "gitlab.com/mgmap/maps"

var ctx *z3.Context
var concreteValuesGlobal *ConcreteValues
var currPathConstrsGlobal *[]z3.Bool
var SymStack *SymbolicStack

func MakeFuzzyInt(name string, a int) int {
	return a
}

func MakeFuzzyBool(name string, b bool) bool {
	return b
}

func initializeGlobals() {
	ctxConfig := z3.NewContextConfig()
	ctxConfig.SetUint("timeout", 5000)
	ctx = z3.NewContext(ctxConfig)

	SymStack = newSymbolicStack()
}

func concolicExecInput(testfunc reflect.Value, cv *ConcreteValues) []reflect.Value {
	// reset global concrete values, path constraints
	concreteValuesGlobal = cv
	newPathConstrs := make([]z3.Bool, 0)
	currPathConstrsGlobal = &newPathConstrs
	SymStack.ClearArgs()

	res := testfunc.Call(make([]reflect.Value, 0))
	return res
}

func concolicForceBranch(branchNum int, branchConds ...z3.Bool) z3.Bool {
	if branchNum < len(branchConds) {
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
		for key, _ := range names.getIntMappings() {
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

func ConcolicExec(testfunc reflect.Value, maxiter int) {
	var hasher maps.Hasher
	hasher = func(o interface{}) uint32 {
		return uint32(o.(z3.Bool).AsAST().Hash())
	}

	var equals maps.Equals
	equals = func(a, b interface{}) bool {
		return a.(z3.Bool).AsAST().Equal(b.(z3.Bool).AsAST())
	}
	seenAlready := maps.NewHashMap(hasher, equals)

	inputs := initialConcreteValueQueue()
	iter := 0
	initializeGlobals()

	for (iter < maxiter) && !(inputs.isEmpty()) {
		iter += 1
		inputThisTime := inputs.dequeue()
		_ = concolicExecInput(testfunc, inputThisTime)
		// fmt.Printf(branchConstrs.AsAST().String())
		for b := 0; b < len(*currPathConstrsGlobal); b++ {
			oppConstr := concolicForceBranch(b, *currPathConstrsGlobal...)
			// fmt.Printf(oppConstr.AsAST().String())
			if seen := seenAlready.Get(oppConstr); seen == nil {
				seenAlready.Put(oppConstr, true)
				newInputFound, newInput := concolicFindInput(oppConstr, inputThisTime)
				if newInputFound {
					newInput.inherit(inputThisTime)
					inputs.enqueue(newInput)
				}
			}
		}
	}
}

func AddPositivePathConstr(constr z3.Bool) {
	*currPathConstrsGlobal = append(*currPathConstrsGlobal, constr)
}

func AddNegativePathConstr(constr z3.Bool) {
	*currPathConstrsGlobal = append(*currPathConstrsGlobal, constr.Not())
}

type Handler struct{}

// an example instrumented function
func rubberducky(iVal int, jVal int) int {
	i := MakeConcolicInt(iVal, SymStack.PopArg().(z3.Int))
	_ = i
	j := MakeConcolicInt(jVal, SymStack.PopArg().(z3.Int))
	_ = j
	SymStack.SetArgsPopped()

	k := i.ConcIntAdd(j)
	if i.ConcIntEq(j).Value {
		AddPositivePathConstr(i.ConcIntEq(j).Z3Expr)
		fmt.Printf("grace is ")
		if i.ConcIntNE(j).Value {
			AddPositivePathConstr(i.ConcIntNE(j).Z3Expr)
			fmt.Println("mean")

			SymStack.PushReturn(k.Z3Expr)
			return k.Value
		} else {
			AddNegativePathConstr(i.ConcIntNE(j).Z3Expr)
			fmt.Printf("pretty")

			fmt.Println(" mean")

			l := i.ConcIntSub(j)
			SymStack.PushReturn(l.Z3Expr)
			return l.Value
		}
	} else {
		AddNegativePathConstr(i.ConcIntEq(j).Z3Expr)
		fmt.Printf("ducks are ")
		if k.ConcIntEq(j).Value {
			AddPositivePathConstr(k.ConcIntEq(j).Z3Expr)
			fmt.Println("great")

			q := k.ConcIntMod(j)
			SymStack.PushReturn(q.Z3Expr)
			return q.Value
		} else {
			AddNegativePathConstr(k.ConcIntEq(j).Z3Expr)
			fmt.Println("cute")

			q := k.ConcIntMul(j)
			SymStack.PushReturn(q.Z3Expr)
			return q.Value
		}
	}
}
