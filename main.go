package main

import (
	"context"
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/debug"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/info"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/list"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/restart"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/start"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd/stop"
	"os"
)

const (
	CLIName = `devenvctl`
)

func main() {
	cmd := rootcmd.New(CLIName)
	cmd.AddCommand(list.Cmd)
	cmd.AddCommand(info.Cmd)
	cmd.AddCommand(start.Cmd)
	cmd.AddCommand(stop.Cmd)
	cmd.AddCommand(restart.Cmd)
	cmd.AddCommand(debug.Cmd)

	if e := cmd.ExecuteContext(context.Background()); e != nil {
		log.New("CLI").Errorf(`Exited with error: %v`, e)
		os.Exit(1)
	}
}
