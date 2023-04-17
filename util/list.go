package util

import (
	"container/list"
	"sync"
)

type PMap struct {
	data *sync.Map
	l    *list.List
	lock *sync.RWMutex
}

func NewPMap() *PMap {
	return &PMap{data: new(sync.Map), l: new(list.List), lock: new(sync.RWMutex)}
}
func (l *PMap) Add(key any, value any) {
	l.l.PushBack(key)
	l.data.Store(key, value)
}
func (l *PMap) IsEmpty() bool {
	return l.l.Len() == 0
}
