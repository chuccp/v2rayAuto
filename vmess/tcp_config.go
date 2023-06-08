package vmess

import (
	core2 "github.com/chuccp/v2rayAuto/core"
)

type TcpConfig struct {
	Id        string
	AlterId   uint32
	CreateNum int
	context   *core2.Context
}

type tcpConfig struct {
	Id   string
	Port int
}
