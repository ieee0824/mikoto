// Package checker implements mikoto's function length and opt-in purity checks.
package checker

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var maxLines int

// Analyzer checks function length and functions marked with //mikoto:pure.
var Analyzer = &analysis.Analyzer{
	Name: "mikoto",
	Doc:  "check function length and opt-in purity rules",
	Run:  run,
}

func init() {
	Analyzer.Flags.IntVar(&maxLines, "max-lines", 80, "maximum number of lines in a function (0 disables)")
}

func run(pass *analysis.Pass) (any, error) {
	pure := map[types.Object]bool{}
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if hasMarker(fn.Doc) {
				pure[pass.TypesInfo.Defs[fn.Name]] = true
			}
		}
	}
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if maxLines > 0 {
				start := pass.Fset.Position(fn.Pos()).Line
				end := pass.Fset.Position(fn.End()).Line
				if n := end - start + 1; n > maxLines {
					pass.Reportf(fn.Pos(), "function %s is %d lines (limit %d)", fn.Name.Name, n, maxLines)
				}
			}
			if pure[pass.TypesInfo.Defs[fn.Name]] {
				checkPure(pass, fn, pure)
			}
		}
	}
	return nil, nil
}

func hasMarker(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, c := range doc.List {
		if strings.TrimSpace(strings.TrimPrefix(c.Text, "//")) == "mikoto:pure" {
			return true
		}
	}
	return false
}

func checkPure(pass *analysis.Pass, fn *ast.FuncDecl, pure map[types.Object]bool) {
	if obj, ok := pass.TypesInfo.Defs[fn.Name].(*types.Func); ok {
		sig := obj.Type().(*types.Signature)
		if !valueTuple(sig.Params()) || !valueTuple(sig.Results()) || (sig.Recv() != nil && !valueOnly(sig.Recv().Type())) {
			pass.Reportf(fn.Name.Pos(), "pure function requires value-only parameters, receiver, and results")
		}
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			pass.Reportf(x.Pos(), "pure function cannot contain a closure")
			return false
		case *ast.GoStmt, *ast.DeferStmt, *ast.SendStmt:
			pass.Reportf(n.Pos(), "pure function cannot start goroutines, defer calls, or send to channels")
			return false
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if _, ok := lhs.(*ast.Ident); !ok {
					pass.Reportf(lhs.Pos(), "pure function cannot assign through a reference")
				}
			}
		case *ast.IncDecStmt:
			if _, ok := x.X.(*ast.Ident); !ok {
				pass.Reportf(x.Pos(), "pure function cannot mutate through a reference")
			}
		case *ast.RangeStmt:
			t := pass.TypesInfo.TypeOf(x.X)
			if t != nil {
				switch t.Underlying().(type) {
				case *types.Map, *types.Chan, *types.Signature:
					pass.Reportf(x.Pos(), "pure function cannot range over a map, channel, or iterator")
				}
			}
		case *ast.UnaryExpr:
			if x.Op == token.ARROW || x.Op == token.AND {
				pass.Reportf(x.Pos(), "pure function cannot receive from a channel or take an address")
			}
		case *ast.StarExpr:
			if !pass.TypesInfo.Types[x].IsType() {
				pass.Reportf(x.Pos(), "pure function cannot dereference a pointer")
			}
		case *ast.Ident:
			if v, ok := pass.TypesInfo.Uses[x].(*types.Var); ok && v.Pkg() != nil && v.Parent() == v.Pkg().Scope() {
				pass.Reportf(x.Pos(), "pure function cannot access package variable %s", x.Name)
			}
		case *ast.CallExpr:
			if pass.TypesInfo.Types[x.Fun].IsType() { // Type conversion.
				break
			}
			var obj types.Object
			switch f := x.Fun.(type) {
			case *ast.Ident:
				obj = pass.TypesInfo.Uses[f]
			case *ast.SelectorExpr:
				obj = pass.TypesInfo.Uses[f.Sel]
			}
			if b, ok := obj.(*types.Builtin); ok {
				switch b.Name() {
				case "len", "cap", "complex", "real", "imag", "min", "max":
					break
				default:
					pass.Reportf(x.Pos(), "pure function cannot call builtin %s", b.Name())
				}
			} else if !pure[obj] {
				pass.Reportf(x.Pos(), "pure function can only call another marked pure function")
			}
		}
		return true
	})
}

func valueTuple(tuple *types.Tuple) bool {
	for i := 0; i < tuple.Len(); i++ {
		if !valueOnly(tuple.At(i).Type()) {
			return false
		}
	}
	return true
}

func valueOnly(t types.Type) bool {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		return t.Kind() != types.UnsafePointer
	case *types.Array:
		return valueOnly(t.Elem())
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			if !valueOnly(t.Field(i).Type()) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
