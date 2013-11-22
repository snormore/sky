package ast

// A slice of Variables.
type Variables []*Variable

func (s Variables) Len() int {
	return len(s)
}

func (s Variables) Less(i, j int) bool {
	if s[i].PropertyId == s[j].PropertyId {
		return s[i].Name < s[j].Name
	}
	return s[i].PropertyId < s[j].PropertyId
}

func (s Variables) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
