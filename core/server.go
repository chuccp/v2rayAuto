package core

type Server interface {
	Start(context *Context) error
	Close() error
	GetKey() string
	GetClient() []string
}
type ApiServer interface {
	Start(context *Context) error
}
