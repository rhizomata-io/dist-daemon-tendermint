package main

import (
	"os"
	"path/filepath"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/tm"
	
	"github.com/tendermint/tendermint/libs/cli"
	
	cmd "github.com/rhizomata-io/dist-daemon-tendermint/cmd/commands"
)

const (
	DefaultBCDir = "chainroot"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.NewInitCmd(),
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		cmd.ShowNodeIDCmd,
		cmd.VersionCmd,
	)
	
	spaces := []string{common.SpaceDaemon, common.SpaceDaemonData}
	
	nodeProvider := tm.NewNodeProvider(spaces)
	// Create & start node
	rootCmd.AddCommand(cmd.NewStartCmd(nodeProvider.NewNode))
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
