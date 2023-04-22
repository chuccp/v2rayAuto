package cert

import "github.com/v2fly/v2ray-core/v5/transport/internet/tls"

type Certificate struct {
	*tls.Certificate
	Domain string
}

func NewCertificate(cert *tls.Certificate, domain string) *Certificate {
	return &Certificate{Certificate: cert, Domain: domain}
}
