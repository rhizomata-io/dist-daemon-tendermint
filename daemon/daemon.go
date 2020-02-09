package daemon

import (
	"time"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/cluster"
	dmCfg "github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/tm/client"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type Daemon struct {
	id             string
	tmNode         *node.Node
	tmCfg          *cfg.Config
	logger         log.Logger
	client         types.Client
	clusterManager *cluster.Manager
	jobManager     *job.Manager
}

func NewDaemon(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config dmCfg.DaemonConfig) (dm *Daemon) {
	dm = &Daemon{tmCfg: tmCfg, logger: logger, tmNode: tmNode}
	dm.client = client.NewClient(tmCfg, logger, types.BasicCdc)
	dm.id = string(dm.tmNode.NodeInfo().ID())
	
	dm.clusterManager = cluster.NewManager(config, logger, dm.client)
	dm.jobManager = job.NewManager(config, logger, dm.client)
	
	return dm
}

func (dm *Daemon) Start() {
	go func() {
		dm.waitReady()
		
		dm.logger.Info("[Dist-Daemon] Start Daemon...", "node_id", dm.tmNode.NodeInfo().ID())
		
		dm.clusterManager.Start()
		dm.jobManager.Start()
	}()
}

func (dm *Daemon) waitReady() {
	threshold := time.Second * 3
	for time.Now().Sub(dm.tmNode.ConsensusState().GetState().LastBlockTime) > threshold {
		time.Sleep(200 * time.Millisecond)
	}
}

func (dm *Daemon) ID() string { return dm.id }

func (dm *Daemon) GetClient() types.Client { return dm.client }

func (dm *Daemon) GetCluster() *cluster.Cluster { return dm.clusterManager.GetCluster() }

func (dm *Daemon) IsLeaderNode() bool { return dm.clusterManager.IsLeaderNode() }

func (dm *Daemon) GetJobRepository() job.Repository { return dm.jobManager.GetRepository() }
