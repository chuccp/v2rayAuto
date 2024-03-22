package vmess

import (
	"encoding/base64"
	c "github.com/chuccp/v2rayAuto/core"
	core "github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	"strconv"
)

type TcpServer struct {
	instance  *core.Instance
	context   *c.Context
	tcpConfig *TcpConfig
}

func (s *TcpServer) Start(context *c.Context) (err error) {
	s.context = context
	s.tcpConfig, err = CreateTcpConfig(context)
	if err != nil {
		return
	}
	err = s.Flush()
	return
}
func (s *TcpServer) Flush() (err error) {
	if s.tcpConfig.CreateNum == 0 {
		return nil
	}
	s.tcpConfig.FlushPort(s.tcpConfig.CreateNum)
	if s.instance != nil {
		err = s.instance.Close()
		if err != nil {
			return
		}
	}
	s.instance, err = CreateTcpServer(s.tcpConfig)
	if err != nil {
		return
	}
	return s.instance.Start()
}

func (s *TcpServer) Close() error {
	return s.instance.Close()
}
func (s *TcpServer) GetKey() string {
	return "tcp"
}
func (s *TcpServer) GetClient() []string {
	urls := make([]string, 0)
	host := s.context.GetCertificate().Domain
	for _, port := range s.tcpConfig.GetShowPorts() {
		name := host + ":" + strconv.Itoa(port)
		config := "{\"v\":\"2\",\"ps\":\"  " +
			name + "   \",\"add\":\"" +
			host + "\",\"port\":\"" +
			strconv.Itoa(port) + "\",\"id\":\"" +
			s.tcpConfig.Id + "\",\"aid\":\"0\",\"scy\":\"auto\",\"net\":\"tcp\",\"type\":\"none\",\"host\":\"" +
			host + "\",\"sni\":\"\",\"alpn\":\"\",\"fp\":\"\",\"allowInsecure\":true}"
		url := "vmess://" + base64.StdEncoding.EncodeToString([]byte(config))
		urls = append(urls, url)
	}
	return urls
}
