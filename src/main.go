package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
// for rewriting
import "reflect"
import "bytes"
import "go/parser"
import "go/ast"
import "go/token"
import "go/printer"
import "golang.org/x/tools/go/ast/astutil"

// import "concolicTypes"

// import "reflect"

// argment is path to example program
var mruNames []string

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

	ast.Print(fset, uninstrumentedAST)
	instrumentedAST := astutil.Apply(uninstrumentedAST, astutil.ApplyFunc(addInstrumentationPre), astutil.ApplyFunc(addInstrumentationPost))

	// concolicExecute(instrumentedAST)
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, instrumentedAST)
	fmt.Println(buf.String())
}

// case *ast.BinaryExpr:
// case *ast.BasicLit:
// if curNode.Kind == token.INT {
// 	// implement replacement
// 	fmt.Printf("quack quack %s\r\n", curNode.Value)
// }

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	return true

}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	fmt.Println(reflect.TypeOf(curNode.Node()))
	switch curNode.Node().(type) {
	case *ast.BasicLit:
		castedNode := curNode.Node().(*ast.BasicLit)
		if castedNode.Kind == token.INT && len(mruNames) > 0 {
			identifier := mruNames[0]
			mruNames = mruNames[1:]

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
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: identifier,
								},
							},
						},
					},
				}
			curNode.Replace(&bruh)
			// TODO implement replacement
		}
		return true
	// case *ast.FuncType:
	case *ast.BlockStmt:
		mruNames = mruNames[1:]
		return true

	case *ast.BinaryExpr:
		// if onlyInts(curNode.(*ast.BinaryExpr)) {

		// }
		// // TODO implement pls kthxbai
		return true
	case *ast.Ident:
		castedNode := curNode.Node().(*ast.Ident)
		typeMetadata := castedNode.Obj
		// if len(castedNode) != 1 {
		// 	panic("oh ducking motherducker")
		// }
		// TODO add to this as we add more types
		if typeMetadata != nil {
			mruNames = append(mruNames, castedNode.Name)
		}
		return true
	default:
		return true
	}
}

func concolicExecute(instrumentedFile ast.Node) {

}
