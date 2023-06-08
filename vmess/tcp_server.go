package vmess

import (
	c "github.com/chuccp/v2rayAuto/core"
	core "github.com/v2fly/v2ray-core/v5"
)

type TcpServer struct {
	instance  *core.Instance
	context   *c.Context
	tcpConfig *TcpConfig
}

func (s *TcpServer) Start(context *c.Context) error {
	return nil
}
func (s *TcpServer) Close() error {
	return nil
}
func (s *TcpServer) Flush() error {
	return nil
}
func (s *TcpServer) GetKey() string {
	return "vmess_tcp"
}
func (s *TcpServer) GetClient() []string {
	return nil
}
