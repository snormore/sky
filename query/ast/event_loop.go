package ast

// EventLoop represents a statement in the query that iterates
// over individual events.
type EventLoop struct {
	Statements Statements
}

func (l *EventLoop) node()      {}
func (l *EventLoop) block()     {}
func (l *EventLoop) statement() {}

// NewEventLoop returns a new EventLoop instance.
func NewEventLoop() *EventLoop {
	return &EventLoop{}
}

// Converts the loop to a string-based representation.
func (l *EventLoop) String() string {
	str := "FOR EACH EVENT\n"
	str += lineStartRegex.ReplaceAllString(l.Statements.String(), "  ") + "\n"
	str += "END"
	return str
}
