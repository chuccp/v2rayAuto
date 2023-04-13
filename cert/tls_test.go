package cert

import (
	"fmt"
	"github.com/chuccp/v2rayAuto/util"
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
func TestSaveCert(t *testing.T) {

	file, err := util.NewFile("key.pem")
	if err != nil {
		t.Log(err)
	}

	exists, _ := file.Exists()
	if !exists {
		cert, err := GetCertificateFromSelf()
		if err != nil {
			log.Fatalln("obtains certificate:", err)
		}
		certPEM, keyPEM := cert.ToPEM()
		err = util.WriteFile("cert.pem", certPEM)
		if err != nil {
			log.Fatalln("WriteFile cert  certificate:", err)
		}
		err = util.WriteFile("key.pem", keyPEM)
		if err != nil {
			log.Fatalln("WriteFile key certificate:", err)
		}
	}
}
func TestCert222(t *testing.T) {

	pem, key, err := LoadCertPem("127.0.0.1", "cooge123@gmail.com", "", 80)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(pem)
	t.Log(key)

}
