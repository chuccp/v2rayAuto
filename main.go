package main

import (
	"github.com/chuccp/v2rayAuto/api"
	"github.com/chuccp/v2rayAuto/v2ray"
	"github.com/gin-gonic/gin"
)

func main() {

	v2ray.Start()

	r := gin.Default()
	r.GET("/d3MuY2oyMDIw.md", api.Subscribe)

	//r.GET("/put", vmess.Api2)
	r.Run(":8082")
}
