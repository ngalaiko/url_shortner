package params

type parseRule func(value []byte) interface{}

func ParseString(value []byte) string {
	return string(value)
}
