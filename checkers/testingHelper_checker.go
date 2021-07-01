package checkers

import (
	"go/ast"
	"go/types"

	"github.com/go-critic/go-critic/checkers/internal/astwalk"
	"github.com/go-critic/go-critic/framework/linter"

	"github.com/go-toolsmith/astfmt"
)

func init() {
	var info linter.CheckerInfo
	info.Name = "testingHelper"
	info.Tags = []string{"style", "experimental"}
	info.Summary = "Detects helper test functions that do not call t.Helper()"
	info.Before = `func myHelper(v int, t *testing.T) {
	if v != 10 {
		t.Fatal("expected 10")
	}
}`
	info.After = `func myHelper(t *testing.T, v int) {
	t.Helper()
	if v != 10 {
		t.Fatal("expected 10")
	}
}`

	collection.AddChecker(&info, func(ctx *linter.CheckerContext) linter.FileWalker {
		return astwalk.WalkerForFuncDecl(&testingHelperChecker{ctx: ctx})
	})
}

type testingHelperChecker struct {
	astwalk.WalkHandler
	ctx *linter.CheckerContext
}

func (c *testingHelperChecker) VisitFuncDecl(decl *ast.FuncDecl) {
	if isUnitTestFunc(c.ctx, decl) {
		return
	}
	if len(decl.Type.Params.List) == 0 {
		return
	}

	typ := c.ctx.TypeOf(decl.Name)
	sig, ok := typ.(*types.Signature)
	if !ok {
		return
	}

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		typ := params.At(i).Type().String()
		if typ != "*testing.T" {
			continue
		}

		if i != 0 {
			c.warnFirstParam(decl)
		}
		if !c.hasFirstCallHelper(decl.Body) {
			c.warnPossibleHelper(decl)
		}
		break
	}
}

func (c *testingHelperChecker) hasFirstCallHelper(body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return true
	}

	expr := body.List[0]
	stmt, ok := expr.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	return astfmt.Sprint(call) == "t.Helper()"
}

func (c *testingHelperChecker) warnFirstParam(fn ast.Node) {
	c.ctx.Warn(fn, "consider to make *testing.T a first parameter")
}

func (c *testingHelperChecker) warnPossibleHelper(fn ast.Node) {
	c.ctx.Warn(fn, "consider to call t.Helper() a first statement")
}
