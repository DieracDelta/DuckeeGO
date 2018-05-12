package main

import (
	// for rewriting
	// "fmt"
	"go/ast"
	// "go/token"
	"golang.org/x/tools/go/ast/astutil"
	// "reflect"
)

// var nodeNumber = 0

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	if _, ok := curNode.Node().(ast.Decl); ok {
		instrumentDeclParentCheckPre(curNode)
	}
	return true
}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	// probably unnecessary to do this twice but it's fine
	if _, ok := curNode.Node().(ast.Decl); ok {
		instrumentDeclParentCheckPre(curNode)
	}
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
	case *ast.CallExpr:
		instrumentCallExpr(curNode)
	case *ast.ReturnStmt:
		instrumentReturnStmt(curNode)
	default:
		// TODO do nothing
	}
	return true
}
