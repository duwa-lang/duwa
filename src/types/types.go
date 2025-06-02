package types

type Type interface {
	String() string
	Equals(Type) bool
}
