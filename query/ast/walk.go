package ast

// Walk traverses the AST in depth-first order.
func Walk(v Visitor, node Node) {
	walk(v, node, nil)
}

func walk(v Visitor, node Node, symtable *Symtable) {
	if node == nil {
		return
	}

	// Generate a new symtable for blocks.
	symtable = NodeSymtable(node, symtable)

	v = v.Visit(node, symtable)
	if v == nil {
		return
	}

	switch n := node.(type) {
	case *Assignment:
		walk(v, n.Target, symtable)
		walk(v, n.Expression, symtable)

	case *BinaryExpression:
		walk(v, n.LHS, symtable)
		walk(v, n.RHS, symtable)

	case *Condition:
		walk(v, n.Expression, symtable)
		walkStatements(v, n.Statements, symtable)

	case *Debug:
		walk(v, n.Expression, symtable)

	case *EventLoop:
		walkStatements(v, n.Statements, symtable)

	case *Field:
		walk(v, n.Expression, symtable)

	case *Query:
		walkVarDecls(v, n.SystemVarDecls, symtable)
		walkVarDecls(v, n.DeclaredVarDecls, symtable)
		walkStatements(v, n.Statements, symtable)

	case *Selection:
		walkFields(v, n.Fields, symtable)

	case *SessionLoop:
		walkStatements(v, n.Statements, symtable)

	case *TemporalLoop:
		walk(v, n.Iterator, symtable)
		walkStatements(v, n.Statements, symtable)

	case *VarDecl:
		if err := symtable.Add(n); err != nil {
			v = v.Error(err)
			if v == nil {
				return
			}
		}
	}

	// Call the visitor on the way out.
	if v, ok := v.(ExitingVisitor); ok {
		v.Exiting(node, symtable)
	}
}

func walkStatements(v Visitor, statements Statements, symtable *Symtable) {
	for _, statement := range statements {
		walk(v, statement, symtable)
	}
}

func walkFields(v Visitor, fields Fields, symtable *Symtable) {
	for _, field := range fields {
		walk(v, field, symtable)
	}
}

func walkVarDecls(v Visitor, decls VarDecls, symtable *Symtable) {
	for _, decl := range decls {
		walk(v, decl, symtable)
	}
}
