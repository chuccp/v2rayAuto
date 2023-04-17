package cert

import (
	"github.com/chuccp/v2rayAuto/util"
	"github.com/v2fly/v2ray-core/v5/common"
	"time"
)

func LoadCertPem(Domain string, Email string, Path string, ValidDay int) ([]byte, []byte, error) {

	certPemFilename := Domain + ".cert.pem"
	keyPemFilename := Domain + ".key.pem"

	file, err := util.NewFile(Path)
	common.Must(err)

	certFile, err := file.Child(certPemFilename)
	common.Must(err)
	keyFile, err := file.Child(keyPemFilename)
	common.Must(err)

	cerExists, err := certFile.Exists()
	common.Must(err)
	keyExists, err := keyFile.Exists()
	common.Must(err)
	if !cerExists || !keyExists {
		return createCertPem(Domain, Email, certFile, keyFile)
	}

	certTime, err := certFile.ModTime()
	common.Must(err)
	keyTime, err := keyFile.ModTime()
	common.Must(err)
	t := time.Now()
	if certTime.Add(time.Hour*time.Duration(ValidDay*24)).Before(t) || keyTime.Add(time.Hour*time.Duration(ValidDay*24)).Before(t) {
		return createCertPem(Domain, Email, certFile, keyFile)
	}
	certPEM, err := certFile.ReadAll()
	common.Must(err)
	keyPEM, err := keyFile.ReadAll()
	common.Must(err)
	return certPEM, keyPEM, nil
}

func createCertPem(Domain string, Email string, certFile *util.File, keyPemFile *util.File) ([]byte, []byte, error) {
	if util.IsIP(Domain) {
		cert, err := GetCertificateFromSelf()
		if err != nil {
			return nil, nil, err
		}
		certPEM, keyPEM := cert.ToPEM()
		err = certFile.WriteBytes(certPEM)
		if err != nil {
			return nil, nil, err
		}
		err = keyPemFile.WriteBytes(keyPEM)
		if err != nil {
			return nil, nil, err
		}
		return certPEM, keyPEM, nil
	} else {
		cert, err := GetCertificateFromLego([]string{Domain}, Email)
		if err != nil {
			return nil, nil, err
		}
		err = certFile.WriteBytes(cert.Certificate)
		if err != nil {
			return nil, nil, err
		}
		err = keyPemFile.WriteBytes(cert.PrivateKey)
		if err != nil {
			return nil, nil, err
		}
		return cert.Certificate, cert.PrivateKey, nil
	}
}
