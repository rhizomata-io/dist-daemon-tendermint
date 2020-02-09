package daemon

import (
	"fmt"
	"time"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/cluster"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/worker"
	
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
	workerManager  *worker.Manager
	jobOrganizer   job.Organizer
}

func NewDaemon(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config common.DaemonConfig) (dm *Daemon) {
	dm = &Daemon{tmCfg: tmCfg, logger: logger, tmNode: tmNode}
	dm.client = client.NewClient(tmCfg, logger, types.BasicCdc)
	dm.id = string(dm.tmNode.NodeInfo().ID())
	
	dm.clusterManager = cluster.NewManager(config, logger, dm.client)
	dm.jobManager = job.NewManager(config, logger, dm.client)
	dm.workerManager = worker.NewManager(config, logger, dm.client)
	return dm
}

func (dm *Daemon) SetJobOrganizer(organizer   job.Organizer){
	dm.jobOrganizer = organizer
}

func (dm *Daemon) Start() {
	go func() {
		dm.waitReady()
		
		dm.logger.Info("[Dist-Daemon] Start Daemon...", "node_id", dm.tmNode.NodeInfo().ID())
		
		dm.clusterManager.Start()
		dm.jobManager.Start()
		dm.workerManager.Start()
		
		if dm.jobOrganizer == nil {
			dm.jobOrganizer = job.NewSimpleOrganizer()
		}
		
		common.SubscribeDaemonEvent(cluster.MemberChangedEventPath,
			"onMemberChanged",
			dm.onMemberChanged)
		
		common.SubscribeDaemonEvent(job.MemberJobsChangedEventPath,
			"onMemberJobsChanged",
			dm.onMemberJobsChanged)
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

// member's jobs changed event handler
func (dm *Daemon) onMemberChanged(event types.Event) {
	dm.logger.Debug(" - [Daemon] onMemberChanged :", event)
	memberChangedEvent := event.(cluster.MemberChangedEvent)
	
	if memberChangedEvent.IsLeader {
		dm.logger.Info("[Daemon] members changed", "members", memberChangedEvent.AliveMemberIDs)
		dm.distributeJobs(memberChangedEvent.AliveMemberIDs)
	}
}

func (dm *Daemon) distributeJobs(aliveMembers []string) {
	allJobs, err := dm.jobManager.GetRepository().GetAllJobs()
	if err != nil {
		if !types.IsNoDataError(err) {
			dm.logger.Error("[Daemon] distributeJobs - GetAllJobs ", err)
		}
		
		return
	}
	membJobMap, err := dm.jobManager.GetRepository().GetAllMemberJobIDs()
	if err != nil {
		dm.logger.Error("[Daemon] distributeJobs - GetAllMemberJobIDs ", err)
		return
	}
	
	newMembJobs, err := dm.jobOrganizer.Distribute(allJobs, aliveMembers, membJobMap)
	
	dm.logger.Info("[Daemon] distributeJobs : ", newMembJobs)
}

// member's jobs changed event handler
func (dm *Daemon) onMemberJobsChanged(event types.Event) {
	dm.logger.Debug(" - [Daemon] onMemberJobsChanged :", event)
	
	memberJobsChangedEvent := event.(job.MemberJobsChangedEvent)
	
	if memberJobsChangedEvent.NodeID != dm.ID() {
		return
	}
	
	dm.logger.Info("[Daemon] member's jobs changed", "nodeID", memberJobsChangedEvent.NodeID,
		"jobs", memberJobsChangedEvent.JobIDs)
	
	jobs, err := dm.jobManager.GetRepository().GetMemberJobs(dm.ID())
	
	if err != nil {
		dm.logger.Error(fmt.Sprintf("[Daemon] cannot get %s's jobs", dm.ID()))
		return
	}
	dm.workerManager.SetJobs(jobs)
}
