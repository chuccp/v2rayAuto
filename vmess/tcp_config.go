package vmess

import (
	"container/list"
	core2 "github.com/chuccp/v2rayAuto/core"
)

type TcpConfig struct {
	Id        string
	AlterId   uint32
	CreateNum int
	ports     *list.List
	context   *core2.Context
}
