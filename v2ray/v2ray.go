package v2ray

import (
	"crypto/x509"
	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/protocol/tls/cert"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	"github.com/v2fly/v2ray-core/v5/proxy/freedom"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess/inbound"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
	"github.com/v2fly/v2ray-core/v5/transport/internet/websocket"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

func Start() {
	userID := "7cc7589e-7ebc-11ec-a352-00163e00bdef"
	serverPort := 8083
	caCert, err := cert.Generate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment|x509.KeyUsageCertSign))
	common.Must(err)
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
						From: uint32(serverPort),
						To:   uint32(serverPort),
					},
					Listen: net.NewIPOrDomain(net.AnyIP),
					StreamSettings: &internet.StreamConfig{
						ProtocolName: "websocket",
						TransportSettings: []*internet.TransportConfig{
							{
								ProtocolName: "websocket",
								Settings: serial.ToTypedMessage(&websocket.Config{
									Path: "/4CErPOyj/",
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
								SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_ZERO},
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
	common.Must(err)
	err = instance.Start()
	common.Must(err)
	log.Println(serverConfig)
}
