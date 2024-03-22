package core

import (
	"github.com/chuccp/v2rayAuto/cert"
	"github.com/chuccp/v2rayAuto/config"
	"github.com/chuccp/v2rayAuto/util"
	"github.com/go-acme/lego/v4/log"
	"github.com/v2fly/VSign/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
	"strconv"
	"strings"
	"sync"
)

type Context struct {
	serverMap   *sync.Map
	config      *config.Config
	certificate *tls.Certificate
	portRange   *net.PortRange
	host        string
	cron        string
	port        int
}

func (v *Context) initConfig() {
	v.host = common.Must2(v.ReadString("core", "url")).(string)
	from, to, err := v.ReadRangeInt("core", "range_port")
	common.Must(err)
	if to > from {
		v.portRange = &net.PortRange{From: uint32(from), To: uint32(to)}
	} else {
		v.portRange = &net.PortRange{To: uint32(from), From: uint32(to)}
	}
	v.cron = common.Must2(v.ReadString("core", "cron")).(string)
	log.Println(v.cron)
}

func (v *Context) GetHost() string {
	return v.host
}
func (v *Context) GetPort() int {
	return v.port
}
func (v *Context) GetNoUsePorts(createNum int) []int {
	return util.GetNoUsePort(int(v.portRange.From), int(v.portRange.To), createNum, []int{v.port})
}

func (v *Context) RegisterServer(server Server) {
	v.serverMap.LoadOrStore(server.GetKey(), server)
}

func (v *Context) GetCertificate() *cert.Certificate {
	return common.Must2(v.createCert()).(*cert.Certificate)
}

func (v *Context) createCert() (*cert.Certificate, error) {
	ports := v.GetNoUsePorts(2)
	domain := common.Must2(v.ReadString("core", "host")).(string)
	email := common.Must2(v.ReadString("core", "email")).(string)
	pem, key, c, k, err := cert.LoadCertPem(domain, email, "", 60, ports[0], ports[1])
	return cert.NewCertificate(&tls.Certificate{Certificate: pem, Key: key, CertificateFile: c, KeyFile: k}, domain), err
}

func (v *Context) ReadString(section string, key string) (string, error) {
	return v.config.ReadString(section, key)
}

func (v *Context) HasSection(section string) bool {
	return v.config.HasSection(section)
}

func (v *Context) ReadInt(section string, key string) (int, error) {
	return v.config.ReadInt(section, key)
}
func (v *Context) ReadRangeInt(section string, key string) (int, int, error) {
	readString, err := v.ReadString(section, key)
	if err != nil {
		return 0, 0, err
	}
	sss := strings.Split(readString, "-")
	a0, err := strconv.Atoi(sss[0])
	if err != nil {
		return 0, 0, err
	}
	a1, err := strconv.Atoi(sss[1])
	if err != nil {
		return 0, 0, err
	}
	return a0, a1, err
}

func (v *Context) RangeServer(f func(server Server)) {
	v.serverMap.Range(func(_, value any) bool {
		f(value.(Server))
		return true
	})
}
