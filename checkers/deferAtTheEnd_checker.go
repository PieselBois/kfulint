package checkers

import (
	"go/ast"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
	"github.com/go-toolsmith/astfmt"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "deferAtTheEnd"
	info.Tags = []string{"diagnostic", "experimental"}
	info.Summary = "Detects calls to defer at the end of function"
	info.Before = `
func() {
	defer os.Remove(filename)
}`
	info.After = `
func() {
	os.Remove(filename)
}`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForFuncDecl(&deferAtTheEnd{ctx: ctx})
	})
}

type deferAtTheEnd struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext
}

func (c *deferAtTheEnd) VisitFuncDecl(funcDecl *ast.FuncDecl) {
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		funDecl, ok := n.(*ast.BlockStmt)
		if ok {
			c.checkDeferBeforeReturn(funDecl)
		}
		return true
	})
}

func (c *deferAtTheEnd) checkDeferBeforeReturn(funcDecl *ast.BlockStmt) {
	retIndex := len(funcDecl.List)
	for i, stmt := range funcDecl.List {
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		if containsCallExpr(retStmt) {
			continue
		}
		retIndex = i
		break

	}
	if retIndex == 0 {
		return
	}

	if deferStmt, ok := funcDecl.List[retIndex-1].(*ast.DeferStmt); ok {
		c.warn(deferStmt)
	}
}

func containsCallExpr(retStmt *ast.ReturnStmt) bool {
	for _, expr := range retStmt.Results {
		if _, ok := expr.(*ast.CallExpr); ok {
			return true
		}
	}
	return false
}

func (c *deferAtTheEnd) warn(deferStmt *ast.DeferStmt) {
	s := astfmt.Sprint(deferStmt)
	if fnlit, ok := deferStmt.Call.Fun.(*ast.FuncLit); ok {
		// To avoid long and multi-line warning messages,
		// collapse the function literals.
		s = "defer " + astfmt.Sprint(fnlit.Type) + "{...}(...)"
	}
	c.ctx.Warn(deferStmt, "%s is placed just before return", s)
}
