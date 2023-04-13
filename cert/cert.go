package cert

import (
	"github.com/chuccp/v2rayAuto/util"
	"time"
)

func LoadCertPem(Domain string, Email string, Path string, ValidDay int) ([]byte, []byte, error) {

	certPemFilename := Domain + ".cert.pem"
	keyPemFilename := Domain + ".key.pem"

	file, err := util.NewFile(Path)
	if err != nil {
		return nil, nil, err
	}

	certFile, err := file.Child(certPemFilename)
	if err != nil {
		return nil, nil, err
	}
	keyFile, err := file.Child(keyPemFilename)
	if err != nil {
		return nil, nil, err
	}

	cerExists, err := certFile.Exists()
	if err != nil {
		return nil, nil, err
	}
	keyExists, err := keyFile.Exists()
	if err != nil {
		return nil, nil, err
	}
	if !cerExists || !keyExists {
		return createCertPem(Domain, Email, certFile, keyFile)
	}

	certTime, err := certFile.ModTime()
	if err != nil {
		return nil, nil, err
	}
	keyTime, err := keyFile.ModTime()
	if err != nil {
		return nil, nil, err
	}
	t := time.Now()
	if certTime.Add(time.Hour*time.Duration(ValidDay*24)).Before(t) || keyTime.Add(time.Hour*time.Duration(ValidDay*24)).Before(t) {
		return createCertPem(Domain, Email, certFile, keyFile)
	}
	certPEM, err := certFile.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	keyPEM, err := keyFile.ReadAll()
	if err != nil {
		return nil, nil, err
	}
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
