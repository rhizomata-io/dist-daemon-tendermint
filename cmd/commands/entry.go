package commands

import (
	"os"
	"path/filepath"
	
	"github.com/tendermint/tendermint/libs/cli"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
)

const (
	DefaultBCDir = "chainroot"
)

func DoCmd(daemonProvider *daemon.BaseProvider){
	rootCmd := InitRootCommand()
	
	AddStartCommand(rootCmd, daemonProvider)
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
