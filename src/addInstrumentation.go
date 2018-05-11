package main

import (
	// for rewriting
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
	// "reflect"
)

var nodeNumber = 0

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	newId := ast.AstId{Id: nodeNumber}
	ast.BinaryExpr
	pointer := curNode.Node().GetId()
	// *pointer = *newId
	nodeNumber++
	// TODO don't really need anything in here yet
	return true
}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	// TODO dead code delete
	// bruh:
	// 	if len(queueOfThings.stage2.parentParent) > 0 {
	// 		if curNode.Node() == queueOfThings.stage2.parentParent[len(queueOfThings.stage2.parentParent)-1] {
	// 			curNode.InsertAfter(queueOfThings.stage2.stmts[len(queueOfThings.stage2.stmts)-1])
	// 			queueOfThings.stage2.Pop(len(queueOfThings.stage2.stmts))
	// 			goto bruh
	// 		}
	// 	}
	exerciseQueueThing(curNode)
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
	case *ast.CallExpr:
		instrumentCallExpr(curNode)
	case *ast.ReturnStmt:
		instrumentReturnStmt(curNode)
	default:
		// TODO do nothing
	}
	return true
}
