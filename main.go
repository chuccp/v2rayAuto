package main

import (
	"github.com/chuccp/v2rayAuto/api"
	"github.com/chuccp/v2rayAuto/config"
	"github.com/chuccp/v2rayAuto/core"
	"github.com/chuccp/v2rayAuto/vmess"
	"github.com/v2fly/v2ray-core/v5/common"
)

func main() {
	readConfig, err := config.ReadConfig("config.ini")
	common.Must(err)
	v2 := core.LoadConfig(readConfig, &api.Server{})
	v2.RegisterServer(&vmess.WsServer{})
	err = v2.Start()
	common.Must(err)
}
