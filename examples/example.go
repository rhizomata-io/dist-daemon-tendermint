package main

import (
	"fmt"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/worker/hello"
	"time"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
)

type ExampleDaemonProvider struct {
	daemon.DefaultProvider
}

func (provider ExampleDaemonProvider) New(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config common.DaemonConfig) *daemon.Daemon {
	provider.DefaultProvider = daemon.DefaultProvider{}
	provider.OnDaemonStarted = provider.onDaemonStarted
	return provider.New(tmCfg, logger, tmNode, config)
}

func (provider ExampleDaemonProvider) onDaemonStarted(dm *daemon.Daemon) {
	
	time.Sleep(4 * time.Second)
	fmt.Println("\n------------- Sample STARTED -------------\n")
	
	//dm.GetJobRepository().RemoveAllJobs()
	
	jobInfo1 := []byte(`{"interval":"0.2s","greet":"hello 0.2s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi1", jobInfo1))
	
	jobInfo2 := []byte(`{"interval":"1.5s","greet":"hello 1.5s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi2", jobInfo2))
	
	jobInfo3 := []byte(`{"interval":"100ms","greet":"hello 100ms" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi3", jobInfo3))
	
	jobInfo4 := []byte(`{"interval":"0.8s","greet":"hi 0.8s" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi4", jobInfo4))
	
	jobInfo5 := []byte(`{"interval":"150ms","greet":"hi 150ms" }`)
	dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi5", jobInfo5))
	
	dm.GetJobRepository().GetJob("hi2")
}
