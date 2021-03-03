package checkers

import (
	"go/ast"

	"github.com/go-critic/go-critic/checkers/internal/astwalk"
	"github.com/go-critic/go-critic/framework/linter"
)

func init() {
	var info linter.CheckerInfo
	info.Name = "randImportChecker"
	info.Tags = []string{"security", "experimental"}
	info.Summary = "Detects when imported package names shadowed in the assignments"
	info.Before = ``
	info.After = ``

	collection.AddChecker(&info, func(ctx *linter.CheckerContext) linter.FileWalker {
		return &randImportChecker{ctx: ctx}
	})
}

type randImportChecker struct {
	astwalk.WalkHandler
	ctx *linter.CheckerContext
}

func (c *randImportChecker) WalkFile(f *ast.File) {
	var cryptoRandImport, mathRandImport *ast.ImportSpec

	for _, importDcl := range f.Imports {
		pkg := importDcl.Path.Value
		println(pkg)
		switch pkg {
		case `"crypto/rand"`:
			cryptoRandImport = importDcl
		case `"math/rand"`:
			mathRandImport = importDcl
		}
	}

	if cryptoRandImport == nil || mathRandImport == nil {
		return
	}
	if cryptoRandImport.Name != nil && mathRandImport.Name == nil {
		c.warn(mathRandImport)
	}
}

func (c *randImportChecker) warn(id ast.Node) {
	c.ctx.Warn(id, "ouch")
}
