package query

type VarRef struct {
	value string
}

func (v *VarRef) String() string {
	return v.value
}
