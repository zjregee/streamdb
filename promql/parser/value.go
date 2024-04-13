package parser

type Value interface {
	Type()   ValueType
	String() string
}

type ValueType string

const (
	ValueTypeNone   ValueType = "none"
	ValueTypeScalar ValueType = "scalar"
	ValueTypeVector ValueType = "vector"
	ValueTypeMatrix ValueType = "matrix"
)

func DocumentedType(vt ValueType) string {
	switch vt {
	case ValueTypeVector:
		return "instant vector"
	case ValueTypeMatrix:
		return "range vector"
	default:
		return string(vt)
	}
}
