package ast

type QueryElement interface {
	Parent() QueryElement
	SetParent(QueryElement)
	Query() *Query
	ElementId() int
	VarRefs() []*VarRef
}

type queryElementImpl struct {
	parent    QueryElement
	elementId int
}

func (e *queryElementImpl) ElementId() int {
	if e.elementId == 0 {
		query := e.Query()
		if query != nil {
			e.elementId = query.NextIdentifier()
		}
	}
	return e.elementId
}

func (e *queryElementImpl) Parent() QueryElement {
	return e.parent
}

func (e *queryElementImpl) SetParent(parent QueryElement) {
	e.parent = parent
}

func (e *queryElementImpl) Query() *Query {
	if e.parent == nil {
		return nil
	}
	return e.parent.Query()
}
