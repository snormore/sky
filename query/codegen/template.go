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
		"metadecl": metadecl,
		"vardecl":  vardecl,
	}
	tmpl = template.New("query").Funcs(m)
	template.Must(tmpl.Parse(string(cursor_tmpl())))
	template.Must(tmpl.Parse(string(event_tmpl())))
	template.Must(tmpl.Parse(string(histogram_tmpl())))
	template.Must(tmpl.Parse(string(average_tmpl())))
	template.Must(tmpl.Parse(string(distinct_tmpl())))
	template.Must(tmpl.Parse(string(util_tmpl())))
	template.Must(tmpl.Parse(string(initialize_tmpl())))
	template.Must(tmpl.Parse(string(aggregate_tmpl())))
	template.Must(tmpl.Parse(string(merge_tmpl())))
	template.Must(tmpl.Parse(string(query_tmpl())))
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

func metadecl(v *query.Variable) string {
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
