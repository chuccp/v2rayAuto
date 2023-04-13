package v2ray

import (
	"encoding/base64"
	c "github.com/chuccp/v2rayAuto/core"
	"github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common/net"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	"strconv"
)

type WsServer struct {
	instance        *core.Instance
	webSocketConfig *WebSocketConfig
	usePorts        []*net.PortRange
}

func (s *WsServer) Start(context *c.Context) error {
	domain, err := context.ReadString("vmess_ws", "domain")
	if err != nil {
		return err
	}
	email, err := context.ReadString("vmess_ws", "email")
	if err != nil {
		return err
	}

	start, end, err := context.ReadRangeInt("vmess_ws", "range_port")
	if err != nil {
		return err
	}
	createNum, err := context.ReadInt("vmess_ws", "create_num")
	if err != nil {
		return err
	}
	s.webSocketConfig = CreateWebSocketConfig(domain, email, uint32(start), uint32(end), uint32(createNum))
	s.instance, s.usePorts, err = CreateWebSocketServer(s.webSocketConfig)
	if err != nil {
		return err
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
	for _, port := range s.usePorts {
		for i := port.FromPort(); i <= port.ToPort(); i++ {
			name := s.webSocketConfig.Domain + "" + strconv.Itoa(int(i))
			config := "{\"v\":\"2\",\"ps\":\"  " +
				name + "   \",\"add\":\"" +
				s.webSocketConfig.Domain + "\",\"port\":\"" +
				strconv.Itoa(int(i)) + "\",\"id\":\"" +
				s.webSocketConfig.Id + "\",\"aid\":\"0\",\"scy\":\"auto\",\"net\":\"ws\",\"type\":\"none\",\"host\":\"" +
				s.webSocketConfig.Domain + "\",\"path\":\"/coke/\",\"tls\":\"tls\",\"sni\":\"\",\"alpn\":\"\",\"fp\":\"\",\"allowInsecure\":true}"
			url := "vmess://" + base64.StdEncoding.EncodeToString([]byte(config))
			urls = append(urls, url)
		}
	}
	return urls
}
