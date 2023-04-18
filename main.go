package main

import (
	"fmt"
	"github.com/chuccp/v2rayAuto/api"
	"github.com/chuccp/v2rayAuto/config"
	"github.com/chuccp/v2rayAuto/core"
	"github.com/chuccp/v2rayAuto/vmess"
	"github.com/mitchellh/go-ps"
	"github.com/v2fly/v2ray-core/v5/common"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
)

func run() {
	readConfig, err := config.ReadConfig("config.ini")
	common.Must(err)
	v2 := core.LoadConfig(readConfig, &api.Server{})
	v2.RegisterServer(&vmess.WsServer{})
	err = v2.Start()
	common.Must(err)
}
func kill(name string) {
	processes, err := ps.Processes()
	if err != nil {
		return
	}
	for _, process := range processes {
		if name == process.Executable() {
			sysType := runtime.GOOS
			if sysType == "linux" {
				cmd := exec.Command("kill", `-9`, strconv.Itoa(process.Pid()))
				cmd.Run()
			}
			if sysType == "windows" {
				cmd := exec.Command("taskkill", `/T`, `/F`, `/pid`, strconv.Itoa(process.Pid()))
				cmd.Run()
			}

		}
	}
}

func main() {
	cmd := "run"
	args := os.Args
	if len(args) > 1 {
		cmd = args[1]
	}
	//dir, _ := os.Getwd()
	fmt.Println("===============================")
	fmt.Println(path.Base(args[0]))

	dir, file := filepath.Split(args[0])
	fmt.Println(dir, "====================", file)
	switch cmd {

	case "start":
		{
			cmd := exec.Command(args[0], "run", "-name-v2rayAuto")
			err := cmd.Start()
			common.Must(err)
		}
	case "run":
		{
			run()
		}
	case "stop":
		{
			kill(file)
		}
	}

}
