package query

// A slice of Variables.
type Variables []*Variable

func (s Variables) Len() int {
	return len(s)
}

func (s Variables) Less(i, j int) bool {
	return s[i].PropertyId < s[j].PropertyId || s[i].Name < s[j].Name
}

func (s Variables) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
