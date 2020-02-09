package job

import (
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	tmevents "github.com/rhizomata-io/dist-daemon-tendermint/tm/events"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

// Manager manager for jobs
type Manager struct {
	nodeID string
	dao    Repository
	logger log.Logger
}

// NewManager ..
func NewManager(config common.DaemonConfig, logger log.Logger, client types.Client) *Manager {
	dao := NewRepository(config, logger, client)
	manager := Manager{nodeID: config.NodeID, dao: dao, logger: logger}
	return &manager
}

func (manager *Manager) GetRepository() Repository {
	return manager.dao
}

func (manager *Manager) Start() {
	jobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathJobs, "")
	tmevents.SubscribeTxEvent(jobsEvtPath, "jobs", func(event tmevents.TxEvent) {
		jobIDs, err := manager.dao.GetAllJobIDs()
		if err != nil {
			manager.logger.Error("[JobManager GetAllJobIDs", err)
		}
		common.PublishDaemonEvent(JobsChangedEvent{
			JobIDs: jobIDs,
		})
	})
	
	memJobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathMemberJobs, manager.nodeID)
	
	tmevents.SubscribeTxEvent(memJobsEvtPath, "mem_jobs", func(event tmevents.TxEvent) {
		nodeID := string(event.Key)
		
		jobIDs, err := manager.dao.GetMemberJobIDs(nodeID)
		if err != nil {
			manager.logger.Error("[JobManager GetAllJobIDs", err)
		}
		common.PublishDaemonEvent(MemberJobsChangedEvent{
			NodeID: nodeID,
			JobIDs: jobIDs,
		})
	})
}
