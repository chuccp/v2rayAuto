package v2ray

import (
	core "github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

type Server struct {
	instance *core.Instance
	key      string
}

func (s *Server) Start() error {
	return s.instance.Start()
}
func (s *Server) Close() error {
	return s.instance.Close()
}
func (s *Server) GetKey() string {
	return s.key
}
