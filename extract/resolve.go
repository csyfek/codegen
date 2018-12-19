package extract

import (
	"fmt"
	"go/ast"
)

func resolveTypeExpression(expr ast.Expr) string {

	switch t := expr.(type) {
	case *ast.StarExpr:
		return "*" + resolveTypeExpression(t.X)
	case *ast.ArrayType:
		return "[]" + resolveTypeExpression(t.Elt)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return resolveTypeExpression(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", resolveTypeExpression(t.Key), resolveTypeExpression(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct"
	default:
		return fmt.Sprintf("_%T_", t)
	}

}
