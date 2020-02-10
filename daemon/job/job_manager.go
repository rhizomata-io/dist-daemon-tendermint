package job

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	tmevents "github.com/rhizomata-io/dist-daemon-tendermint/tm/events"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

// Manager manager for jobs
type Manager struct {
	common.Context
	nodeID string
	dao    Repository
	//logger log.Logger
}

// NewManager ..
func NewManager(context common.Context) *Manager {
	dao := NewRepository(context.GetConfig(), context, context.GetClient())
	manager := Manager{Context: context, dao: dao, nodeID: context.GetNodeID()}
	return &manager
}

func (manager *Manager) GetRepository() Repository {
	return manager.dao
}

func (manager *Manager) Start() {
	
	var jobsChanged bool
	
	jobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathJobs, "")
	tmevents.SubscribeTxEvent(jobsEvtPath, "jobs", func(event tmevents.TxEvent) {
		jobsChanged = true
	})
	
	
	tmevents.SubscribeBlockEvent(tmevents.CommitEventPath, "check-jobs-changed", func(event types.Event) {
		if jobsChanged {
			manager.Info("[JobManager] Jobs changed.")
			jobIDs, err := manager.dao.GetAllJobIDs()
			if err != nil {
				manager.Error("[JobManager] Jobs changed", err)
				return
			}
			
			manager.Info("[JobManager] Jobs changed.", jobIDs)
			
			common.PublishDaemonEvent(JobsChangedEvent{
				JobIDs: jobIDs,
			})
			
			jobsChanged = false
		}
	})
	
	
	memJobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathMemberJobs, manager.nodeID)
	
	tmevents.SubscribeTxEvent(memJobsEvtPath, "mem_jobs", func(event tmevents.TxEvent) {
		nodeID := string(event.Key)
		
		jobIDs, err := manager.dao.GetMemberJobIDs(nodeID)
		if err != nil {
			manager.Error("[JobManager GetAllJobIDs", err)
		}
		common.PublishDaemonEvent(MemberJobsChangedEvent{
			NodeID: nodeID,
			JobIDs: jobIDs,
		})
	})
}
