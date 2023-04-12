package v2ray

import "sync"

var serverMap = new(sync.Map)

func RegisterServer(server *Server) {
	serverMap.LoadOrStore(server.GetKey(), server)
}
func RangeServer(f func(server *Server)) {
	serverMap.Range(func(_, value any) bool {
		f(value.(*Server))
		return true
	})
}
func GetServer() []*Server {
	servers := make([]*Server, 0)
	serverMap.Range(func(_, value any) bool {
		servers = append(servers, value.(*Server))
		return true
	})
	return servers
}
