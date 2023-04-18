package vmess

import (
	"encoding/base64"
	c "github.com/chuccp/v2rayAuto/core"
	"github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	"strconv"
)

type WsServer struct {
	instance        *core.Instance
	webSocketConfig *WebSocketConfig
	context         *c.Context
}

func (s *WsServer) Start(context *c.Context) (err error) {
	s.context = context
	s.webSocketConfig, err = CreateWebSocketConfig(context)
	if err != nil {
		return
	}
	err = s.Flush()
	return
}
func (s *WsServer) Flush() (err error) {
	s.webSocketConfig.flushPort()
	if s.instance != nil {
		err = s.instance.Close()
		if err != nil {
			return
		}
	}
	s.instance, err = CreateWebSocketServer(s.webSocketConfig)
	if err != nil {
		return
	}
	return s.instance.Start()
}

func (s *WsServer) Close() error {
	return s.instance.Close()
}
func (s *WsServer) GetKey() string {
	return "WebSocketOverTls"
}
func (s *WsServer) GetClient() []string {
	urls := make([]string, 0)
	host := s.context.GetHost()
	for _, port := range s.webSocketConfig.getPorts() {
		name := host + ":" + strconv.Itoa(port)
		config := "{\"v\":\"2\",\"ps\":\"  " +
			name + "   \",\"add\":\"" +
			host + "\",\"port\":\"" +
			strconv.Itoa(port) + "\",\"id\":\"" +
			s.webSocketConfig.Id + "\",\"aid\":\"0\",\"scy\":\"auto\",\"net\":\"ws\",\"type\":\"none\",\"host\":\"" +
			host + "\",\"path\":\"" + "/coke_" + strconv.Itoa(port) + "/" + "\",\"tls\":\"tls\",\"sni\":\"\",\"alpn\":\"\",\"fp\":\"\",\"allowInsecure\":true}"
		url := "vmess://" + base64.StdEncoding.EncodeToString([]byte(config))
		urls = append(urls, url)
	}
	return urls
}
