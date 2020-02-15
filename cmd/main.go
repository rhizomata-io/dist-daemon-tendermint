package main

import (
	cmd "github.com/rhizomata-io/dist-daemon-tendermint/cmd/commands"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
)

func main() {
	daemonProvider := &daemon.BaseProvider{}
	cmd.DoCmd(daemonProvider)
}
