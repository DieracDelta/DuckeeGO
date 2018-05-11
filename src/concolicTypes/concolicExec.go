package concolicTypes

import "fmt"
import "reflect"
import "github.com/aclements/go-z3/z3"
import "gitlab.com/mgmap/maps"

var ctx *z3.Context

func MakeFuzzyInt(name string, a int) {

}

func MakeFuzzyBool(name string, b bool) {

}

func setGlobalContext() {
	ctxConfig := z3.NewContextConfig()
	ctxConfig.SetUint("timeout", 5000)
	ctx = z3.NewContext(ctxConfig)
}

func concolicExecInput(testfunc reflect.Value, concreteValues *ConcreteValues) ([]reflect.Value, *[]z3.Bool) {
	var currPathConstrs []z3.Bool 
	args := []reflect.Value{reflect.ValueOf(concreteValues), reflect.ValueOf(&currPathConstrs)}
	res := testfunc.Call(args)
	return res, &currPathConstrs
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

func concolicExec(testfunc reflect.Value, maxiter int) {
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
	setGlobalContext()

	for (iter < maxiter) && !(inputs.isEmpty()) {
		iter += 1
		inputThisTime := inputs.dequeue()
		_, branchConstrs := concolicExecInput(testfunc, inputThisTime)
		// fmt.Printf(branchConstrs.AsAST().String())
		for b := 0; b < len(*branchConstrs); b++ {
			oppConstr := concolicForceBranch(b, *branchConstrs...)
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

func addPositivePathConstr(currPathConstrs *[]z3.Bool, constr ConcolicBool) {
	*currPathConstrs = append(*currPathConstrs, constr.z3Expr)
}

func addNegativePathConstr(currPathConstrs *[]z3.Bool, constr ConcolicBool) {
	*currPathConstrs = append(*currPathConstrs, constr.z3Expr.Not())
}

type Handler struct{}

func (h Handler) Rubberducky(cv *ConcreteValues, currPathConstrs *[]z3.Bool) int {
	var i ConcolicInt
	var j ConcolicInt
	i = MakeConcolicIntVar(cv, "i")
	j = MakeConcolicIntVar(cv, "j")
	k := i.ConcIntAdd(j)
	b := i.ConcIntEq(j)
	if b.Value {
		addPositivePathConstr(currPathConstrs, b)
		fmt.Printf("grace is ")
		b1 := i.ConcIntNE(j)
		if b1.Value {
			addPositivePathConstr(currPathConstrs, b1)
			fmt.Printf("mean")
		} else {
			addNegativePathConstr(currPathConstrs, b1)
			fmt.Printf("very helpful")
		}
	} else {
		addNegativePathConstr(currPathConstrs, b)
		fmt.Printf("ducks ")
		b1 := k.ConcIntEq(j)
		if b1.Value {
			addPositivePathConstr(currPathConstrs, b1)
			fmt.Printf("are great")
		} else {
			addNegativePathConstr(currPathConstrs, b1)
			fmt.Printf("are cute")
		}
	}
	fmt.Println()

	var x ConcolicInt
	var y ConcolicInt
	x = MakeConcolicIntVar(cv, "x")
	y = MakeConcolicIntVar(cv, "y")
	b2 := x.ConcIntGE(y)
	if b2.Value {
		addPositivePathConstr(currPathConstrs, b2)
		fmt.Printf("grace ")
		b3 := x.ConcIntLT(y)
		if b3.Value {
			addPositivePathConstr(currPathConstrs, b3)
			fmt.Printf("< ")
		} else {
			addNegativePathConstr(currPathConstrs, b3)
			fmt.Printf("> ")
		}
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
