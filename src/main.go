package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
import "go/parser"
import "go/ast"
import "go/token"

// import "reflect"

// argment is path to example program
func main() {
	fmt.Print("mr duck\r\n")
	fset := token.NewFileSet()
	// TODO add more files  by including more args
	fileName := os.Args[1]

	parsedFile, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	// ast.Print(fset, parsedFile)
	addInstrumentation(fset, parsedFile)

	concolic_execute(parsedFile)
}

// DFS it, recursively
// TODO switch to iterative cuz recursion sux
func addInstrumentation(fset *token.FileSet, parsedFile *ast.File) {
	queue := parsedFile.Decls
	for _, curNode := range queue {
		instrumentNode(curNode)
	}
}

func instrumentNode(curNode interface{}) {
	switch curNode := curNode.(type) {
	// case *ast.GenDecl:
	// 	switch curNode.Tok {
	// 	case token.CONST:
	// 	case token.TYPE:
	// 	case token.VAR:
	// 	case token.IMPORT:
	// 	}
	// 	fmt.Print(curNode)
	case *ast.FuncDecl:
		instrumentNode(curNode.Body)
		// fmt.Print(curNode)
	case *ast.BlockStmt:
		for _, ele := range curNode.List {
			instrumentNode(ele)
		}
	case *ast.DeclStmt:
		instrumentNode(curNode.Decl)
	case *ast.GenDecl:
		switch curNode.Tok {
		case token.VAR:
			for _, ele := range curNode.Specs {
				instrumentNode(ele)
			}
		}
	case *ast.ValueSpec:
		for _, aNode := range curNode.Values {
			instrumentNode(aNode)
		}
	case *ast.BinaryExpr:
		// TODO implement pls kthxbai
	case *ast.BasicLit:
		if curNode.Kind == token.INT {
			fmt.Printf("quack quack %s\r\n", curNode.Value)
		}
	}

	// fmt.Print(curNode)
	// fmt.Print("\r\n\r\n\r\n")
}

func concolic_execute(instrumenetedFile *ast.File) {
}
