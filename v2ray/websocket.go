package v2ray

import (
	"github.com/chuccp/v2rayAuto/cert"
	"github.com/chuccp/v2rayAuto/util"
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
)

type WebSocketConfig struct {
	Path           string
	FromPort       uint32
	ToPort         uint32
	Id             string
	AlterId        uint32
	CamouflageHost string
	Domain         string
	CreateNum      uint32
	Email          string
	ramPort        *ramPort
}

type ramPort struct {
	fromPort    uint32
	createNum   uint32
	toPort      uint32
	perFromPort uint32
}

func CreateWebSocketConfig(domain string, email string, FromPort uint32, ToPort uint32, createNum uint32) *WebSocketConfig {
	uuid := uuid.New()
	return &WebSocketConfig{
		FromPort: FromPort,
		ToPort:   ToPort,
		Path:     "/coke/",
		AlterId:  0,
		Id:       uuid.String(),
		Domain:   domain,
		Email:    email,
		ramPort:  &ramPort{fromPort: FromPort, createNum: createNum, toPort: 0},
	}
}
func (wsc *WebSocketConfig) getRamPort() *ramPort {
	rp := wsc.ramPort
	if rp.toPort == 0 {
		rp.perFromPort = rp.fromPort
		rp.toPort = rp.perFromPort + (rp.createNum - 1)
	} else {
		rp.fromPort = rp.perFromPort
		rp.perFromPort = rp.fromPort + rp.createNum
		rp.toPort = rp.perFromPort + (rp.createNum - 1)
	}
	if rp.toPort > wsc.ToPort {
		rp.fromPort = wsc.FromPort
		rp.toPort = rp.fromPort + (rp.createNum - 1)
	}
	return wsc.ramPort
}
func (wsc *WebSocketConfig) getPortRanges() []*net.PortRange {
	ramPort := wsc.getRamPort()
	pr := &net.PortRange{From: ramPort.perFromPort, To: ramPort.toPort}
	portRanges := util.GetNoUsePort(pr)
	if ramPort.fromPort < ramPort.perFromPort {
		portRanges = append(portRanges, &net.PortRange{From: ramPort.fromPort, To: ramPort.perFromPort - 1})
	}
	return portRanges
}
func (wsc *WebSocketConfig) toWsConfig(portRange *net.PortRange) *wsConfig {
	return &wsConfig{Path: wsc.Path, FromPort: portRange.From, ToPort: portRange.To, Id: wsc.Id}
}

type wsConfig struct {
	Path     string
	FromPort uint32
	ToPort   uint32
	Key      string
	Id       string
}

func getWebSocketInboundHandlerConfigs(webSocketConfig *WebSocketConfig, pem []byte, key []byte) ([]*core.InboundHandlerConfig, []*net.PortRange, error) {
	inboundHandlerConfigs := make([]*core.InboundHandlerConfig, 0)
	portRanges := webSocketConfig.getPortRanges()
	for _, portRange := range portRanges {
		InboundHandlerConfig, err := getWebSocketInboundHandlerConfig(webSocketConfig.toWsConfig(portRange), pem, key)
		if err != nil {
			return nil, portRanges, err
		}
		inboundHandlerConfigs = append(inboundHandlerConfigs, InboundHandlerConfig)
	}
	return inboundHandlerConfigs, portRanges, nil
}
func getWebSocketInboundHandlerConfig(webSocketConfig *wsConfig, pem []byte, key []byte) (*core.InboundHandlerConfig, error) {
	userID := webSocketConfig.Id
	inboundHandlerConfig := &core.InboundHandlerConfig{
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortRange: &net.PortRange{
				From: webSocketConfig.FromPort,
				To:   webSocketConfig.ToPort,
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
						Certificate: []*tls.Certificate{{
							Certificate: pem,
							Key:         key,
							Usage:       tls.Certificate_AUTHORITY_ISSUE,
						}},
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

func CreateWebSocketServer(webSocketConfig *WebSocketConfig) (*core.Instance, []*net.PortRange, error) {

	pem, key, err := cert.LoadCertPem(webSocketConfig.Domain, webSocketConfig.Email, "", 80)
	if err != nil {
		return nil, nil, err
	}

	inboundHandlerConfigs, portRanges, err := getWebSocketInboundHandlerConfigs(webSocketConfig, pem, key)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}
	return instance, portRanges, nil
}
