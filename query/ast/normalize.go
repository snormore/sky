package ast

// Normalize walks the AST tree and fixes nodes that are incomplete. The
// main use for this is to insert event loops to ensure forward motion of the
// execution.
func Normalize(q *Query) {
	ForEach(q, normalize)
}

func normalize(n Node) {
	switch n := n.(type) {
	case *Query:
		normalizeQuery(n)
	}
}

// normalizeQuery ensures that at least one loop exists in the query.
func normalizeQuery(n *Query) {
	if containsLoop(n.Statements) {
		return
	}

	// Separate declarations from executable statements.
	preamble, body := make(Statements, 0), make(Statements, 0)
	for _, statement := range n.Statements {
		switch statement.(type) {
		case *VarDecl:
			preamble = append(preamble, statement)
		default:
			body = append(body, statement)
		}
	}

	// Insert body as an event loop after the preamble.
	preamble = append(preamble, &EventLoop{Statements: body})
	n.Statements = preamble
}

// containsLoop returns true if a set of statements contains at least one loop.
func containsLoop(statements Statements) bool {
	for _, statement := range statements {
		switch statement.(type) {
		case *EventLoop, *SessionLoop, *TemporalLoop:
			return true
		}
	}
	return false
}
