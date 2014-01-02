package ast

// Walk traverses the AST in depth-first order.
func Walk(v Visitor, node Node) {
	walk(v, node, nil)
}

func walk(v Visitor, node Node, symtable *Symtable) Visitor {
	if node == nil {
		return v
	}

	// Generate a new symtable for blocks.
	symtable = NodeSymtable(node, symtable)

	v = v.Visit(node, symtable)
	if v == nil {
		return nil
	}

	switch n := node.(type) {
	case *Assignment:
		if v = walk(v, n.Target, symtable); v == nil {
			return nil
		}
		if v = walk(v, n.Expression, symtable); v == nil {
			return nil
		}

	case *BinaryExpression:
		if v = walk(v, n.LHS, symtable); v == nil {
			return nil
		}
		if v = walk(v, n.RHS, symtable); v == nil {
			return nil
		}

	case *Condition:
		if v = walk(v, n.Expression, symtable); v == nil {
			return nil
		}
		if v = walkStatements(v, n.Statements, symtable); v == nil {
			return nil
		}

	case *Debug:
		if v = walk(v, n.Expression, symtable); v == nil {
			return nil
		}

	case *EventLoop:
		if v = walkStatements(v, n.Statements, symtable); v == nil {
			return nil
		}

	case *Field:
		if v = walk(v, n.Expression, symtable); v == nil {
			return nil
		}

	case *Query:
		if v = walkVarDecls(v, n.SystemVarDecls, symtable); v == nil {
			return nil
		}
		if v = walkVarDecls(v, n.dynamicVarDecls, symtable); v == nil {
			return nil
		}
		if v = walkVarDecls(v, n.DeclaredVarDecls, symtable); v == nil {
			return nil
		}
		if v = walkStatements(v, n.Statements, symtable); v == nil {
			return nil
		}

	case *Selection:
		if v = walkVarRefs(v, n.Dimensions, symtable); v == nil {
			return nil
		}
		if v = walkFields(v, n.Fields, symtable); v == nil {
			return nil
		}

	case *SessionLoop:
		if v = walkStatements(v, n.Statements, symtable); v == nil {
			return nil
		}

	case *TemporalLoop:
		if v = walk(v, n.Iterator, symtable); v == nil {
			return nil
		}
		if v = walkStatements(v, n.Statements, symtable); v == nil {
			return nil
		}

	case *VarDecl:
		if err := symtable.Add(n); err != nil {
			if v = v.Error(err); v == nil {
				return nil
			}
		}
	}

	// Call the visitor on the way out.
	if v, ok := v.(ExitingVisitor); ok {
		v.Exiting(node, symtable)
	}

	return v
}

func walkStatements(v Visitor, statements Statements, symtable *Symtable) Visitor {
	for _, statement := range statements {
		if v = walk(v, statement, symtable); v == nil {
			return nil
		}
	}
	return v
}

func walkFields(v Visitor, fields Fields, symtable *Symtable) Visitor {
	for _, field := range fields {
		if v = walk(v, field, symtable); v == nil {
			return nil
		}
	}
	return v
}

func walkVarDecls(v Visitor, decls VarDecls, symtable *Symtable) Visitor {
	for _, decl := range decls {
		if v = walk(v, decl, symtable); v == nil {
			return v
		}
	}
	return v
}

func walkVarRefs(v Visitor, refs []*VarRef, symtable *Symtable) Visitor {
	for _, ref := range refs {
		if v = walk(v, ref, symtable); v == nil {
			return v
		}
	}
	return v
}
