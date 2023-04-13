package v2ray

type Server interface {
	Start() error
	Close() error
	GetKey() string
	GetClient() []string
}
