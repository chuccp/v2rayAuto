package db

import (
	"github.com/mitchellh/go-ps"
	"testing"
)

func TestName(t *testing.T) {
	processes, err := ps.Processes()
	if err != nil {
		return
	}
	for i, process := range processes {
		t.Log(i, process.Executable())
	}
}
