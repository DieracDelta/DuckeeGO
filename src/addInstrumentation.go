package main

import (
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
)

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	if _, ok := curNode.Node().(ast.Decl); ok {
		instrumentDeclParentCheckPre(curNode)
	}

	switch curNode.Node().(type) {
	case *ast.BasicLit:
		instrumentBasicLitPre(curNode)
	case *ast.AssignStmt:
		instrumentAssignStmtPre(curNode)
	case *ast.CompositeLit:
		instrumentCompositeLitPre(curNode)
	case *ast.FuncDecl:
		instrumentFuncDeclPre(curNode)
	case *ast.IndexExpr:
		instrumentIndexExprPre(curNode)
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
		// do nothing
	}
	return true
}
