package ast

type BooleanLiteral struct {
	queryElementImpl
	value bool
}

func (l *BooleanLiteral) VarRefs() []*VarRef {
	return []*VarRef{}
}

func (l *BooleanLiteral) Codegen() (string, error) {
	return l.String(), nil
}

func (l *BooleanLiteral) String() string {
	if l.value {
		return "true"
	} else {
		return "false"
	}
}
