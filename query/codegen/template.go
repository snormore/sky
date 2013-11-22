package codegen

import (
	"fmt"
	"text/template"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/tmpl/bin"
)

// The template used for generating Lua code from queries.
var tmpl *template.Template

func init() {
	m := template.FuncMap{
		"metadecl": metadecl,
		"vardecl":  vardecl,
	}
	tmpl = template.New("query").Funcs(m)
	template.Must(tmpl.Parse(bin.Tmpl("/cursor.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/event.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/histogram.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/average.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/distinct.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/util.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/initialize.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/aggregate.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/merge.tmpl")))
	template.Must(tmpl.Parse(bin.Tmpl("/query.tmpl")))
}

// vardecl returns the C field definition for a variable.
func vardecl(v *ast.Variable) string {
	if v.IsSystemVariable() || v.Name == "timestamp" {
		return ""
	}

	var ctype string
	switch v.DataType {
	case core.StringDataType:
		ctype = "sky_string_t"
	case core.FactorDataType, core.IntegerDataType:
		ctype = "int32_t"
	case core.FloatDataType:
		ctype = "double"
	case core.BooleanDataType:
		ctype = "bool"
	default:
		return ""
	}

	return fmt.Sprintf("%s _%s;", ctype, v.Name)
}

func metadecl(v *ast.Variable) string {
	if v.IsSystemVariable() {
		return ""
	}
	switch v.DataType {
	case core.StringDataType:
		return fmt.Sprintf("%v = function(event) return ffi.string(event._%v.data, event._%v.length) end,", v.Name, v.Name, v.Name)
	case core.FactorDataType, core.IntegerDataType, core.FloatDataType, core.BooleanDataType:
		return fmt.Sprintf("%v = function(event) return event._%v end,", v.Name, v.Name)
	}
	return ""
}
