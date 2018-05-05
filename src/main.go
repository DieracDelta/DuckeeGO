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

// import "concolicTypes"

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
			bruh :=
				ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "concolicTypes",
						},
						Sel: &ast.Ident{
							Name: "ConcolicInt",
						},
					},
					Elts: []ast.Expr{
						castedNode,
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "symTypes",
								},
								Sel: &ast.Ident{
									Name: "SymInt",
								},
							},
							Elts: []ast.Expr{
								castedNode,
							},
						},
					},
				}
			return &bruh, true
			// implement replacement
		}
		return curNode, true
	case *ast.BinaryExpr:
		// if onlyInts(curNode.(*ast.BinaryExpr)) {

		// }
		// // TODO implement pls kthxbai
		return curNode, true
	default:
		return curNode, true
	}
}
