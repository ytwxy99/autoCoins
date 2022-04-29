package policy

type Policy interface {
	Target(args ...interface{}) interface{}
}
