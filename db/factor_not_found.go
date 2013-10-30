package db

func NewFactorNotFound(text string) error {
	return &FactorNotFound{text}
}

type FactorNotFound struct {
	s string
}

func (e *FactorNotFound) Error() string {
	return e.s
}
