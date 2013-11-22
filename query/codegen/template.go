package codegen

import (
	"fmt"
	"text/template"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query"
)

// The template used for generating Lua code from queries.
var tmpl *template.Template

func init() {
	m := template.FuncMap{
		"vardecl": vardecl,
	}
	tmpl = template.Must(template.New("query").Funcs(m).Parse(string(query_tmpl())))
	template.Must(tmpl.Parse(string(event_decl_tmpl())))
}

// vardecl returns the C field definition for a variable.
func vardecl(v *query.Variable) string {
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
