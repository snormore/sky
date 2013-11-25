package ast

type Fields []*Field

// AggregateFields returns a list of aggregate fields in this collection.
func (s Fields) AggregateFields() bool {
	var ret Fields
	for _, f := range s {
		if f.IsAggregate() {
			ret = append(ret, f)
		}
	}
	return ret
}

// NonAggregateFields returns a list of non-aggregate fields in this collection.
func (s Fields) NonAggregateFields() bool {
	var ret Fields
	for _, f := range s {
		if !f.IsAggregate() {
			ret = append(ret, f)
		}
	}
	return ret
}
