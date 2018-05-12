package main

import (
	// for rewriting
	// "fmt"
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
	// "reflect"
)

// var nodeNumber = 0

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	switch curNode.Node().(type) {
	case *ast.BasicLit:
		instrumentBasicLitPre(curNode)
	case *ast.AssignStmt:
		instrumentAssignStmtPre(curNode)
	case *ast.CompositeLit:
		instrumentCompositeLitPre(curNode)
	}
	return true
}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	switch curNode.Node().(type) {
	// the idea is to find a binary expression
	// then check if it contains an int type (or function that returns int type)
	// replace with the node with callexpr if it does
	case *ast.BinaryExpr:
		instrumentBinaryExprPost(curNode)
	case *ast.UnaryExpr:
		instrumentUnaryExprPost(curNode)
	case *ast.BasicLit:
		instrumentBasicLitPost(curNode)
	case *ast.CompositeLit:
		instrumentCompositeLitPost(curNode)
	case *ast.AssignStmt:
		instrumentAssignStmtPost(curNode)
	case *ast.IncDecStmt:
		instrumentIncDecStmtPost(curNode)
	// case *ast.BlockStmt:
	case *ast.Ident:
		instrumentIdentPost(curNode)
	case *ast.IfStmt:
		instrumentIfStmtPost(curNode)
	case *ast.FuncDecl:
		instrumentFuncDeclPost(curNode)
	case *ast.CallExpr:
		instrumentCallExprPost(curNode)
	case *ast.ReturnStmt:
		instrumentReturnStmtPost(curNode)
	default:
		// TODO do nothing
	}
	return true
}
