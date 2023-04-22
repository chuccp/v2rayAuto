package vmess

import (
	c "github.com/chuccp/v2rayAuto/core"
	core "github.com/v2fly/v2ray-core/v5"
)

type TcpServer struct {
	instance *core.Instance
	context  *c.Context
}
