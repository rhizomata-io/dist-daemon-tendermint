package daemon

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/cluster"
	dmCfg "github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/tm/client"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type Daemon struct {
	id             string
	tmNode         *node.Node
	tmCfg         *cfg.Config
	logger         log.Logger
	client         types.Client
	clusterManager *cluster.Manager
}

func NewDaemon(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config dmCfg.DaemonConfig) (dm *Daemon) {
	dm = &Daemon{tmCfg: tmCfg, logger: logger, tmNode: tmNode}
	dm.client = client.NewClient(tmCfg, logger, types.BasicCdc)
	dm.id = string(dm.tmNode.NodeInfo().ID())
	dm.clusterManager = cluster.NewManager(config, logger,dm.client)
	
	return dm
}

func (dm *Daemon) ID() string { return dm.id }

func (dm *Daemon) CetClient() types.Client { return dm.client }

func (dm *Daemon) Start() {
	dm.logger.Info("Start Daemon...", "node_id", dm.tmNode.NodeInfo().ID())
	dm.clusterManager.Start()
}
