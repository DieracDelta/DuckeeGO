package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
// for rewriting
import "bytes"
import "github.com/fatih/astrewrite"
import "go/parser"
import "go/ast"
import "go/token"
import "go/printer"

// import "reflect"

// argment is path to example program
func main() {
	fmt.Print("mr duck\r\n")
	fset := token.NewFileSet()
	// TODO add more files  by including more args
	filePath := os.Args[1]

	uninstrumentedAST, err := parser.ParseFile(fset, filePath, nil, 0)

	if err != nil {
		panic(err)
	}

	ast.Print(fset, uninstrumentedAST)
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
	switch curNode.(type) {
	case *ast.BasicLit:
		if curNode.(*ast.BasicLit).Kind == token.INT {
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
