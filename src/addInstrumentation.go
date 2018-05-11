package main

import (
	// for rewriting
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
	// "reflect"
)

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	// TODO don't really need anything in here yet
	return true
}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	// fmt.Println(reflect.TypeOf(curNode.Node()))
	switch curNode.Node().(type) {
	// the idea is to find a binary expression
	// then check if it contains an int type (or function that returns int type)
	// replace with the node with callexpr if it does
	// case *ast.
	case *ast.BinaryExpr:
		instrumentBinaryExpr(curNode)
	case *ast.UnaryExpr:
		instrumentUnaryExpr(curNode)
	case *ast.BasicLit:
		instrumentBasicLit(curNode)
	case *ast.AssignStmt:
		instrumentAssignStmt(curNode)
	case *ast.IncDecStmt:
		instrumentIncDecStmt(curNode)
	// case *ast.BlockStmt:
	case *ast.Ident:
		instrumentIdent(curNode)
	case *ast.IfStmt:
		instrumentIfStmt(curNode)
	case *ast.FuncDecl:
		instrumentFuncDecl(curNode)
	default:
		// TODO do nothing
	}
	return true
}
