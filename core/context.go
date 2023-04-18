package core

import (
	"github.com/chuccp/v2rayAuto/cert"
	"github.com/chuccp/v2rayAuto/config"
	"github.com/chuccp/v2rayAuto/util"
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
}

func (v *Context) GetHost() string {
	if len(v.host) > 0 {
		return v.host
	}
	v.host = common.Must2(v.ReadString("core", "host")).(string)
	return v.host
}
func (v *Context) GetPortRange() *net.PortRange {
	if v.portRange != nil {
		return v.portRange
	}
	from, to, err := v.ReadRangeInt("core", "range_port")
	common.Must(err)
	if to > from {
		v.portRange = &net.PortRange{From: uint32(from), To: uint32(to)}
	} else {
		v.portRange = &net.PortRange{To: uint32(from), From: uint32(to)}
	}
	return v.portRange
}
func (v *Context) RegisterServer(server Server) {
	v.serverMap.LoadOrStore(server.GetKey(), server)
}

func (v *Context) GetCertificate() *tls.Certificate {
	return common.Must2(v.createCert()).(*tls.Certificate)
}

func (v *Context) createCert() (*tls.Certificate, error) {
	pr := v.GetPortRange()
	ports := util.GetNoUsePort(int(pr.From), int(pr.To), 2, []int{})
	domain := common.Must2(v.ReadString("tls", "domain"))
	email := common.Must2(v.ReadString("tls", "email"))
	pem, key, err := cert.LoadCertPem(domain.(string), email.(string), "", 80, ports[0], ports[1])
	return &tls.Certificate{Certificate: pem, Key: key}, err
}

func (v *Context) ReadString(section string, key string) (string, error) {
	return v.config.ReadString(section, key)
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
