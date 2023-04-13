package v2ray

import (
	"encoding/base64"
	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common/net"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	"strconv"
)

type WsServer struct {
	instance        *core.Instance
	webSocketConfig *WebSocketConfig
	usePorts        []*net.PortRange
}

func (s *WsServer) Start() error {

	return s.instance.Start()
}
func (s *WsServer) Close() error {
	return s.instance.Close()
}
func (s *WsServer) GetKey() string {
	return s.webSocketConfig.Key
}
func (s *WsServer) GetClient() []string {
	urls := make([]string, 0)
	for _, port := range s.usePorts {
		for i := port.FromPort(); i <= port.ToPort(); i++ {
			name := s.webSocketConfig.Host + "" + strconv.Itoa(int(i))
			config := "{\"v\":\"2\",\"ps\":\"  " +
				name + "   \",\"add\":\"" +
				s.webSocketConfig.Host + "\",\"port\":\"" +
				strconv.Itoa(int(i)) + "\",\"id\":\"" +
				s.webSocketConfig.Id + "\",\"aid\":\"0\",\"scy\":\"auto\",\"net\":\"ws\",\"type\":\"none\",\"host\":\"" +
				s.webSocketConfig.Host + "\",\"path\":\"/coke/\",\"tls\":\"tls\",\"sni\":\"\",\"alpn\":\"\",\"fp\":\"\",\"allowInsecure\":true}"
			url := "vmess://" + base64.StdEncoding.EncodeToString([]byte(config))
			urls = append(urls, url)
		}
	}
	return urls
}
