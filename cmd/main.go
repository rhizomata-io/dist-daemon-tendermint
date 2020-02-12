package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/worker/hello"
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
	
	daemonProvider := daemon.DefaultProvider{
		OnDaemonStarted: onDaemonStarted,
	}
	
	// Create & start node
	rootCmd.AddCommand(cmd.NewStartCmd(nodeProvider.NewNode, daemonProvider.New))
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func onDaemonStarted(dm *daemon.Daemon) {
	
	time.Sleep(4 * time.Second)
	fmt.Println("\n------------- Sample STARTED -------------\n")
	
	//dm.GetJobRepository().RemoveAllJobs()
	
	jobInfo1 := []byte(`{"interval":"0.2s","greet":"hello 0.2s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi1",jobInfo1))
	
	jobInfo2 := []byte(`{"interval":"1.5s","greet":"hello 1.5s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi2",jobInfo2))
	
	jobInfo3 := []byte(`{"interval":"100ms","greet":"hello 100ms" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi3", jobInfo3))
	
	jobInfo4 := []byte(`{"interval":"0.8s","greet":"hi 0.8s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi4", jobInfo4))
	
	jobInfo5 := []byte(`{"interval":"150ms","greet":"hi 150ms" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi5", jobInfo5))
	
	dm.GetJobRepository().GetJob("hi2")
}