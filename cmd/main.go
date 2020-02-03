package main

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/tm"
	"os"
	"path/filepath"
	
	cmd "github.com/rhizomata-io/dist-daemon-tendermint/cmd/commands"
	"github.com/tendermint/tendermint/libs/cli"
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
	
	// Create & start node
	rootCmd.AddCommand(cmd.NewStartCmd(tm.NewNode))
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
