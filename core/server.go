package core

type Server interface {
	Start(context *Context) error
	Close() error
	Flush() error
	GetKey() string
	GetClient() []string
}
type ApiServer interface {
	Start(context *Context) error
}
