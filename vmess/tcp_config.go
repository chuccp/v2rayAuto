package vmess

import (
	core2 "github.com/chuccp/v2rayAuto/core"
	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/common/uuid"
	"github.com/v2fly/v2ray-core/v5/proxy/freedom"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess/inbound"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

type TcpConfig struct {
	*core2.PortConfig
	Id        string
	AlterId   uint32
	CreateNum int
	context   *core2.Context
}
type tcpConfig struct {
	Id   string
	Port int
}

func CreateTcpConfig(context *core2.Context) (*TcpConfig, error) {
	createNum, err := context.ReadInt("vmess_tcp", "create_num")
	log.Println("CreateTcpConfig", createNum)
	if err != nil {
		return nil, err
	}
	uuid := uuid.New()
	return &TcpConfig{
		AlterId:    0,
		Id:         uuid.String(),
		CreateNum:  createNum,
		context:    context,
		PortConfig: core2.NewPortConfig(context),
	}, nil
}

func (ws *TcpConfig) getTcpInboundHandlerConfigs(tcpConfig *TcpConfig) ([]*core.InboundHandlerConfig, error) {
	inboundHandlerConfigs := make([]*core.InboundHandlerConfig, 0)
	for _, port := range ws.GetPorts() {
		wss := tcpConfig.toTCPConfig(port)
		InboundHandlerConfig, err := ws.getTcpConfigInboundHandlerConfig(wss)
		if err != nil {
			return nil, err
		}
		inboundHandlerConfigs = append(inboundHandlerConfigs, InboundHandlerConfig)
	}
	return inboundHandlerConfigs, nil
}

func (ws *TcpConfig) toTCPConfig(port int) *tcpConfig {
	return &tcpConfig{Id: ws.Id, Port: port}
}

func (ws *TcpConfig) getTcpConfigInboundHandlerConfig(tcpConfig *tcpConfig) (*core.InboundHandlerConfig, error) {
	userID := tcpConfig.Id
	inboundHandlerConfig := &core.InboundHandlerConfig{

		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortRange: &net.PortRange{
				From: uint32(tcpConfig.Port),
				To:   uint32(tcpConfig.Port),
			},
			Listen: net.NewIPOrDomain(net.AnyIP),
			//StreamSettings: &internet.StreamConfig{
			//	ProtocolName: "http",
			//	TransportSettings: []*internet.TransportConfig{
			//		{
			//			ProtocolName: "http",
			//			Settings:     serial.ToTypedMessage(&http.Config{}),
			//		},
			//	},
			//},
		}),

		ProxySettings: serial.ToTypedMessage(&inbound.Config{
			User: []*protocol.User{
				{
					Account: serial.ToTypedMessage(&vmess.Account{
						Id:               userID,
						AlterId:          0,
						SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_CHACHA20_POLY1305},
					}),
				},
			},
		}),

		//ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
		//	PortRange: &net.PortRange{
		//		From: uint32(tcpConfig.Port),
		//		To:   uint32(tcpConfig.Port),
		//	},
		//	Listen: net.NewIPOrDomain(net.AnyIP),
		//	StreamSettings: &internet.StreamConfig{
		//		ProtocolName: "tcp",
		//		TransportSettings: []*internet.TransportConfig{
		//			{
		//				ProtocolName: "tcp",
		//				Settings:     serial.ToTypedMessage(&tcp.Config{}),
		//			},
		//		},
		//		//SecurityType: serial.GetMessageType(&tcp.Config{}),
		//		//SecuritySettings: []*anypb.Any{
		//		//	serial.ToTypedMessage(&tcp.Config{}),
		//		//},
		//	}}),
		//ProxySettings: serial.ToTypedMessage(&inbound.Config{
		//	User: []*protocol.User{
		//		{
		//			Account: serial.ToTypedMessage(&vmess.Account{
		//				Id:               userID,
		//				SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_CHACHA20_POLY1305},
		//			}),
		//		},
		//	},
		//}),
	}
	return inboundHandlerConfig, nil
}
func CreateTcpServer(tcpConfig *TcpConfig) (*core.Instance, error) {
	inboundHandlerConfigs, err := tcpConfig.getTcpInboundHandlerConfigs(tcpConfig)
	if err != nil {
		return nil, err
	}
	serverConfig := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Inbound: inboundHandlerConfigs,
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}
	instance, err := core.New(serverConfig)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
