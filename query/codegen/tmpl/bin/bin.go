package bin

// Tmpl returns the template contents by path name.
func Tmpl(path string) string {
	if fn, ok := go_bindata[path]; ok {
		return string(fn())
	}
	return ""
}
