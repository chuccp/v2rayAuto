package core

import (
	"github.com/chuccp/v2rayAuto/config"
	"github.com/robfig/cron/v3"
	"github.com/v2fly/v2ray-core/v5/common"
	"log"
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

func (v *V2ray) createCert() {
	log.Println("证书生成")
	v.context.createCert()
}

func (v *V2ray) Start() error {
	v.context.initConfig()
	v.RangeServer(func(server Server) {
		log.Println("Start", server.GetKey())
		err := server.Start(v.context)
		common.Must(err)
	})
	go func() {
		err := v.apiServer.Start(v.context)
		common.Must(err)
	}()
	cr := cron.New(cron.WithSeconds())
	cr.AddFunc(v.context.cron, func() {
		log.Println("==================重启服务")
		v.RangeServer(func(server Server) {
			err := server.Flush()
			common.Must(err)
		})
		go func() {
			err := v.apiServer.Start(v.context)
			common.Must(err)
		}()
	})
	cr.Start()
	return nil
}
func New(config *config.Config, apiServer ApiServer) *V2ray {
	v := &V2ray{apiServer: apiServer, config: config, context: &Context{config: config, serverMap: new(sync.Map)}}
	return v
}
