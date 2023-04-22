package core

import "container/list"

type PortConfig struct {
	showPorts []int
	ports     *list.List
	context   *Context
}

func NewPortConfig(context *Context) *PortConfig {
	return &PortConfig{showPorts: make([]int, 0), context: context, ports: new(list.List)}

}
func (ws *PortConfig) GetShowPorts() []int {
	return ws.showPorts
}
func (ws *PortConfig) GetPorts() []int {
	showPorts := make([]int, 0)
	for ele := ws.ports.Front(); ele != nil; ele = ele.Next() {
		port := ele.Value.(int)
		showPorts = append(showPorts, port)
	}
	return showPorts
}

func (ws *PortConfig) FlushPort(createNum int) {
	if ws.ports.Len() == 0 {
		ws.showPorts = ws.context.GetNoUsePorts(createNum)
		for _, port := range ws.showPorts {
			ws.ports.PushBack(port)
		}
	} else {
		ws.showPorts = ws.context.GetNoUsePorts(createNum)
		for _, port := range ws.showPorts {
			ws.ports.PushBack(port)
		}
	}
	for {
		if ws.ports.Len() <= (createNum * 2) {
			break
		} else {
			ws.ports.Remove(ws.ports.Front())
		}
	}
}
