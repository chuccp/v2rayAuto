package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/v2fly/v2ray-core/v5/common/protocol/tls/cert"
)

type CertUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *CertUser) GetEmail() string {
	return u.Email
}
func (u CertUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *CertUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func GetCertificateFromLego(domain []string, email string) (*certificate.Resource, error) {
	// 创建myUser用户对象。新对象需要email和私钥才能启动，私钥我们自己生成
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	myUser := CertUser{
		Email: email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// 这里配置密钥的类型和密钥申请的地址，记得上线后替换成 lego.LEDirectoryProduction ，测试环境下就用 lego.LEDirectoryStaging
	config.CADirURL = lego.LEDirectoryProduction
	config.Certificate.KeyType = certcrypto.RSA2048

	// 创建一个client与CA服务器通信
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 这里需要挑战我们申请的证书，必须监听80和443端口，这样才能让Let's Encrypt访问到咱们的服务器
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	if err != nil {
		return nil, err
	}
	err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443"))
	if err != nil {
		return nil, err
	}

	// 把这个客户端注册，传递给myUser用户里
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: domain, // 这里如果有多个，就写多个就好了，可以是多个域名
		Bundle:  true,   // 这里如果是true，将把颁发者证书一起返回，也就是返回里面certificates.IssuerCertificate
	}
	// 开始申请证书
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	// 申请完了后，里面会带有证书的PrivateKey Certificate，都为[]byte格式，需要存储的自行转为string即可
	return certificates, nil
}
func GetCertificateFromSelf() (*cert.Certificate, error) {
	return cert.Generate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment|x509.KeyUsageCertSign))
}
