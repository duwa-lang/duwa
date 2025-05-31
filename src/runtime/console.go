package runtime

type Console interface {
	Read() (string, error)
	Clear() error
}
