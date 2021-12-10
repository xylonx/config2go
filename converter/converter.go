package converter

import "reflect"

type Node struct {
	// inner node properties
	// the tag key of the node
	TagKey string
	// the struct field name
	FieldName string
	// the node type: string, bool, int64, float64, Slice, Map and so on
	Type reflect.Kind

	// sort with alphabeta
	Child []Node

	Error error
}

type Converter interface {
	ConvertToNodeTree() (*Node, error)
}
