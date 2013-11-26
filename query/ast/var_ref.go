package ast

// VarRef represents a reference to a variable in the current scope in a query.
type VarRef struct {
	Name string
}

func (v *VarRef) node()       {}
func (v *VarRef) expression() {}

// NewVarRef creates a new VarRef instance.
func NewVarRef() *VarRef {
	return &VarRef{}
}

func (v *VarRef) String() string {
	return "@" + v.Name
}
