package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
// for rewriting
//import "reflect"
import "bytes"
import "github.com/fatih/astrewrite"
import "go/parser"
import "go/ast"
import "go/token"
import "go/printer"

// import "reflect"

// argment is path to example program
func main() {
	if false {
		fmt.Print("mr duck\r\n")
	}
	fset := token.NewFileSet()
	// TODO add more files  by including more args
	filePath := os.Args[1]

	uninstrumentedAST, err := parser.ParseFile(fset, filePath, nil, 0)

	if err != nil {
		panic(err)
	}

	// ast.Print(fset, uninstrumentedAST)
	instrumentedAST := astrewrite.Walk(uninstrumentedAST, addInstrumentation)

	// concolicExecute(instrumentedAST)
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, instrumentedAST)
	fmt.Println(buf.String())
}

func concolicExecute(instrumentedFile ast.Node) {

}

// case *ast.BinaryExpr:
// case *ast.BasicLit:
// if curNode.Kind == token.INT {
// 	// implement replacement
// 	fmt.Printf("quack quack %s\r\n", curNode.Value)
// }

func addInstrumentation(curNode ast.Node) (ast.Node, bool) {
	// fmt.Println(reflect.TypeOf(curNode))
	switch curNode.(type) {
	case *ast.BasicLit:
		castedNode := curNode.(*ast.BasicLit)
		if castedNode.Kind == token.INT {
			innardExpr :=
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: "SymInt",
					},
					Args: []ast.Expr{castedNode},
				}
			arguments := []ast.Expr{
				castedNode,
				innardExpr,
			}
			bruh :=
				ast.CallExpr{
					Fun: &ast.Ident{
						Name: "ConcolicInt",
					},
					Args:     arguments,
					Ellipsis: token.NoPos}
			return &bruh, true
			// implement replacement
		}
		return curNode, true
	case *ast.BinaryExpr:
		// // TODO implement pls kthxbai
		return curNode, true
	default:
		return curNode, true
	}
}
