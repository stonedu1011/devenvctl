package main

import (
	"context"
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/info"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/list"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/start"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/stop"
	"os"
)

const (
	CLIName = `devenvctl`
	BuildVersion = `unknown`
)

func main() {
	cmd := rootcmd.New(CLIName, BuildVersion)
	cmd.AddCommand(list.Cmd)
	cmd.AddCommand(info.Cmd)
	cmd.AddCommand(start.Cmd)
	cmd.AddCommand(stop.Cmd)
	cmd.AddCommand(stop.Cmd)

	if e := cmd.ExecuteContext(context.Background()); e != nil {
		log.New("CLI").Errorf(`Exited with error: %v`, e)
		os.Exit(1)
	}
}
