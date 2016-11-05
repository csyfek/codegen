package extractor

import (
	"fmt"
	"go/ast"
)

func getExprType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return "*" + getExprType(t.X)
	case *ast.ArrayType:
		return "[]" + getExprType(t.Elt)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return getExprType(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("Unknown Type: %T", t)
	}
}
