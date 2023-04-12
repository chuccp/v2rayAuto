package v2ray

import (
	"crypto/x509"
	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/protocol/tls/cert"
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
	Host           string
	CreateNum      uint32
	ramPort        *ramPort
	Key            string
}

type ramPort struct {
	fromPort  uint32
	createNum uint32
	toPort    uint32
}

func CreateWebSocketConfig(host string, FromPort uint32, ToPort uint32, createNum uint32, key string) *WebSocketConfig {
	uuid := uuid.New()
	return &WebSocketConfig{
		FromPort: FromPort,
		ToPort:   ToPort,
		Path:     "/coke/",
		AlterId:  0,
		Id:       uuid.String(),
		Host:     host,
		ramPort:  &ramPort{fromPort: FromPort, createNum: createNum, toPort: 0},
		Key:      key,
	}
}
func (wsc *WebSocketConfig) getRamPort() *ramPort {
	rp := wsc.ramPort
	if rp.toPort == 0 {
		rp.toPort = rp.fromPort + (rp.createNum - 1)
	} else {
		if rp.createNum > 1 {
			rp.fromPort = rp.toPort - (rp.createNum / 2) + 1
			rp.toPort = rp.fromPort + (rp.createNum - 1)
		} else {
			rp.fromPort = rp.fromPort + 1
			rp.toPort = rp.fromPort
		}
	}
	if rp.toPort > wsc.ToPort {
		rp.fromPort = wsc.FromPort
		rp.toPort = rp.fromPort + (rp.createNum - 1)
	}
	return wsc.ramPort
}
func (wsc *WebSocketConfig) GetClientUrl() string {

	return ""
}

func CreateWebSocketServer(webSocketConfig *WebSocketConfig) (*Server, error) {
	ramPort := webSocketConfig.getRamPort()
	userID := webSocketConfig.Id
	caCert, err := cert.Generate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment|x509.KeyUsageCertSign))
	if err != nil {
		return nil, err
	}
	certPEM, keyPEM := caCert.ToPEM()
	serverConfig := &core.Config{
		App: []*anypb.Any{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: &net.PortRange{
						From: ramPort.fromPort,
						To:   ramPort.toPort,
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
									Certificate: certPEM,
									Key:         keyPEM,
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
			},
		},
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
	return &Server{instance: instance, key: webSocketConfig.Key}, nil

}
