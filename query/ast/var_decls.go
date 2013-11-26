package ast

// A slice of VarDecl objects.
type VarDecls []*VarDecl

func (s VarDecls) Len() int {
	return len(s)
}

func (s VarDecls) Less(i, j int) bool {
	if s[i].Id == s[j].Id {
		return s[i].Name < s[j].Name
	}
	return s[i].Id < s[j].Id
}

func (s VarDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
