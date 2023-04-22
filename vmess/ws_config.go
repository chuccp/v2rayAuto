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
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
	"github.com/v2fly/v2ray-core/v5/transport/internet/websocket"
	"google.golang.org/protobuf/types/known/anypb"
	"strconv"
)

type WebSocketConfig struct {
	*core2.PortConfig
	Path      string
	Id        string
	AlterId   uint32
	CreateNum int
	context   *core2.Context
}

func CreateWebSocketConfig(context *core2.Context) (*WebSocketConfig, error) {
	createNum, err := context.ReadInt("vmess_ws", "create_num")
	if err != nil {
		return nil, err
	}
	uuid := uuid.New()
	return &WebSocketConfig{
		AlterId:    0,
		Id:         uuid.String(),
		CreateNum:  createNum,
		context:    context,
		PortConfig: core2.NewPortConfig(context),
	}, nil
}

type wsConfig struct {
	Path string
	Id   string
	Port int
}

func (ws *WebSocketConfig) toWsConfig(port int) *wsConfig {
	return &wsConfig{Path: "/coke_" + strconv.Itoa(port) + "/", Id: ws.Id, Port: port}
}

func (ws *WebSocketConfig) getWebSocketInboundHandlerConfigs(webSocketConfig *WebSocketConfig) ([]*core.InboundHandlerConfig, error) {
	inboundHandlerConfigs := make([]*core.InboundHandlerConfig, 0)
	for _, port := range ws.GetPorts() {
		wss := webSocketConfig.toWsConfig(port)
		InboundHandlerConfig, err := ws.getWebSocketInboundHandlerConfig(wss)
		if err != nil {
			return nil, err
		}
		inboundHandlerConfigs = append(inboundHandlerConfigs, InboundHandlerConfig)
	}
	return inboundHandlerConfigs, nil
}
func (ws *WebSocketConfig) getWebSocketInboundHandlerConfig(webSocketConfig *wsConfig) (*core.InboundHandlerConfig, error) {
	userID := webSocketConfig.Id
	inboundHandlerConfig := &core.InboundHandlerConfig{
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortRange: &net.PortRange{
				From: uint32(webSocketConfig.Port),
				To:   uint32(webSocketConfig.Port),
			},
			Listen: net.NewIPOrDomain(net.AnyIP),
			StreamSettings: &internet.StreamConfig{
				ProtocolName: "websocket",
				TransportSettings: []*internet.TransportConfig{
					{
						ProtocolName: "websocket",
						Settings: serial.ToTypedMessage(&websocket.Config{
							Path: webSocketConfig.Path,
						}),
					},
				},
				SecurityType: serial.GetMessageType(&tls.Config{}),
				SecuritySettings: []*anypb.Any{
					serial.ToTypedMessage(&tls.Config{
						Certificate: []*tls.Certificate{ws.context.GetCertificate().Certificate},
					}),
				},
			}}),
		ProxySettings: serial.ToTypedMessage(&inbound.Config{
			User: []*protocol.User{
				{
					Account: serial.ToTypedMessage(&vmess.Account{
						Id:               userID,
						SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AES128_GCM},
					}),
				},
			},
		}),
	}
	return inboundHandlerConfig, nil
}

func CreateWebSocketServer(webSocketConfig *WebSocketConfig) (*core.Instance, error) {
	inboundHandlerConfigs, err := webSocketConfig.getWebSocketInboundHandlerConfigs(webSocketConfig)
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
