package daemon

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/worker"
	"github.com/rhizomata-io/dist-daemon-tendermint/tm"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
)


type Provider func(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon


type BaseProvider struct {
	Factories       []worker.Factory
	OnDaemonStarted func(*Daemon)
}



func (provider BaseProvider) NewDaemon(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon {
	dm := NewDaemon(tmCfg, logger, tmNode, config, daemonApp)
	for _, fac := range provider.Factories {
		dm.RegisterWorkerFactory(fac)
	}
	
	common.SubscribeDaemonEvent(StartedEventPath, "onDaemonStarted", func(event types.Event) {
		if provider.OnDaemonStarted != nil {
			provider.OnDaemonStarted(dm)
		}
	})
	return dm
}
