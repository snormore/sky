package ast

// Walk traverses the AST in depth-first order.
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	v = v.Visit(node)
	if v == nil {
		return
	}

	switch n := node.(type) {
	case *Assignment:
		Walk(v, n.Target)
		Walk(v, n.Expression)

	case *BinaryExpression:
		Walk(v, n.LHS)
		Walk(v, n.RHS)

	case *Condition:
		Walk(v, n.Expression)
		walkStatements(v, n.Statements)

	case *Debug:
		Walk(v, n.Expression)

	case *EventLoop:
		walkStatements(v, n.Statements)

	case *Field:
		Walk(v, n.Expression)

	case *Query:
		walkVarDecls(v, n.SystemVarDecls)
		walkVarDecls(v, n.DeclaredVarDecls)
		walkStatements(v, n.Statements)

	case *Selection:
		walkFields(v, n.Fields)

	case *SessionLoop:
		walkStatements(v, n.Statements)

	case *TemporalLoop:
		Walk(v, n.Iterator)
		walkStatements(v, n.Statements)
	}
}

func walkStatements(v Visitor, statements Statements) {
	for _, statement := range statements {
		Walk(v, statement)
	}
}

func walkFields(v Visitor, fields Fields) {
	for _, field := range fields {
		Walk(v, field)
	}
}

func walkVarDecls(v Visitor, decls VarDecls) {
	for _, decl := range decls {
		Walk(v, decl)
	}
}
