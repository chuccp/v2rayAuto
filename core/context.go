package core

import (
	"github.com/chuccp/v2rayAuto/config"
	"strconv"
	"strings"
	"sync"
)

type Context struct {
	serverMap *sync.Map
	config    *config.Config
}

func (v *Context) RegisterServer(server Server) {
	v.serverMap.LoadOrStore(server.GetKey(), server)
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
