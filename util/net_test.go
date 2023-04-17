package util

import (
	"testing"
)

func TestName(t *testing.T) {

}

func TestRegx(t *testing.T) {

	pps := GetNoUsePort(1, 8000, 1000, []int{11})
	for _, pp := range pps {

		t.Log(pp)
	}

}
