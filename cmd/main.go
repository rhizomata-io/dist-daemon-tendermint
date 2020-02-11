package main

import (
	"encoding/json"
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
	
	dm.GetJobRepository().RemoveAllJobs()
	
	//jobInfo1, _ := json.Marshal(`{"interval":"2s","greet":"hello 2s" }`)
	//dm.GetJobRepository().PutJobIfNotExist(job.New(hello.FactoryName, jobInfo1))
	//
	//jobInfo2, _ := json.Marshal(`{"interval":"1.5s","greet":"hello 1.5s" }`)
	//dm.GetJobRepository().PutJobIfNotExist(job.New(hello.FactoryName, jobInfo2))
	//
	//jobInfo3, _ := json.Marshal(`{"interval":"400ms","greet":"hello 400ms" }`)
	//dm.GetJobRepository().PutJobIfNotExist(job.New(hello.FactoryName, jobInfo3))
	
	jobInfo4, _ := json.Marshal(`{"interval":"3s","greet":"hi 3s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi1", jobInfo4))
	
	jobInfo5, _ := json.Marshal(`{"interval":"1.7s","greet":"hi 1.7s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi2", jobInfo5))
	
	jobInfo6, _ := json.Marshal(`{"interval":"500ms","greet":"hi 500ms" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi3", jobInfo6))
	
	//dm.GetJobRepository().GetJob("hi2")
	//dm.GetJobRepository().GetAllJobs()
}
