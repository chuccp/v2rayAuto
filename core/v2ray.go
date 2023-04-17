package core

import (
	"github.com/chuccp/v2rayAuto/config"
	"github.com/robfig/cron/v3"
	"github.com/v2fly/v2ray-core/v5/common"
	"sync"
)

type V2ray struct {
	config    *config.Config
	apiServer ApiServer
	context   *Context
}

func (v *V2ray) RegisterServer(server Server) {
	v.context.RegisterServer(server)
}
func (v *V2ray) RangeServer(f func(server Server)) {
	v.context.RangeServer(f)
}
func (v *V2ray) Start() error {
	v.RangeServer(func(server Server) {
		err := server.Start(v.context)
		common.Must(err)
	})
	cr := cron.New(cron.WithSeconds())
	cr.AddFunc("0 0 0 * * *", func() {
		v.RangeServer(func(server Server) {
			err := server.Flush()
			common.Must(err)
		})
	})
	cr.Start()
	v.apiServer.Start(v.context)
	return nil
}
func LoadConfig(config *config.Config, apiServer ApiServer) *V2ray {
	v := &V2ray{apiServer: apiServer, config: config, context: &Context{config: config, serverMap: new(sync.Map)}}
	return v
}
