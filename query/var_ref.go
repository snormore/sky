package query

type VarRef struct {
	queryElementImpl
	value string
}

func (v *VarRef) String() string {
	return v.value
}
