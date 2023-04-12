package tls

import (
	"fmt"
	"github.com/v2fly/v2ray-core/v5/common/protocol/tls/cert"
	"log"
	"testing"
)

func TestCert(t *testing.T) {
	cs, err := GetCertificateFromLego([]string{"ws.cooge.top"}, "cooge123@gmail.com")
	if err != nil {
		log.Fatalln("obtains certificate:", err)
	}
	cert, err := cert.ParseCertificate(cs.Certificate, cs.PrivateKey)
	if err != nil {

	}
	certPEM, keyPEM := cert.ToPEM()
	fmt.Println(certPEM, keyPEM)
}
