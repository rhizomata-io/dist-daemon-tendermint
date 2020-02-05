package cluster

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

type Manager struct {
	config              config.DaemonConfig
	logger              log.Logger
	cluster             *Cluster
	dao                 *DAO
	running             bool
	memberChangeHandler func(aliveMembers []string)
}

func NewManager(config config.DaemonConfig, logger log.Logger, client types.Client) *Manager {
	cluster := newCluster(config.ChainID)
	dao := DAO{config:config, logger: logger, client: client}
	
	manager := &Manager{
		config:              config,
		logger:              logger,
		cluster:             cluster,
		dao:                 &dao,
	}
	
	localMemb := Member{
		NodeID:    config.NodeID,
		Name:      config.NodeName,
		heartbeat: time.Now(),
		leader:    false,
		alive:     true,
		local:     true,
	}
	
	cluster.putMember(&localMemb)
	cluster.localMember = &localMemb
	
	return manager
}

func (manager *Manager) Start() {
	manager.running = true
	
	err := manager.dao.PutMember(manager.cluster.localMember)
	
	if err != nil {
		manager.logger.Error("Cannot PutMember.", err)
		panic(err)
	}
	
	err = manager.dao.PutHeartbeat(manager.cluster.localMember.NodeID)
	
	if err != nil {
		manager.logger.Error("Cannot send Heartbeat.", err)
		panic(err)
	}
	
	go func() {
		for manager.running {
			time.Sleep(time.Duration(manager.config.HeartbeatInterval))
			err := manager.dao.PutHeartbeat(manager.config.NodeID)
			if err != nil {
				manager.logger.Error("Cannot send Heartbeat.", err)
			}
		}
	}()
	
	go func() {
		for manager.running {
			err := manager.dao.GetHeartbeats(manager.handleHeartbeat)
			if err != nil {
				manager.logger.Error("[FATAL] Cannot check heartbeats.", err)
			}
			manager.checkLeader()
			time.Sleep(time.Duration(manager.config.CheckHeartbeatInterval))
		}
	}()
	
	manager.logger.Info("[INFO-Cluster] Start Cluster Manager.")
}

func (manager *Manager)handleHeartbeat(nodeid string, time time.Time) {

	
}


func (manager *Manager)checkLeader() {

}