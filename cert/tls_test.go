package cert

import (
	"github.com/chuccp/v2rayAuto/util"
	"log"
	"testing"
)

func TestCert(t *testing.T) {

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

}
