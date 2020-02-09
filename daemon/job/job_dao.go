package job

import (
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathJobs       = "jobs"
	PathMemberJobs = "membjobs"
)

// jobDao kv store model for job
type jobDao struct {
	config common.DaemonConfig
	logger log.Logger
	client types.Client
}

var _ Repository = (*jobDao)(nil)

func NewRepository(config common.DaemonConfig, logger log.Logger, client types.Client) Repository {
	dao := &jobDao{config: config, logger: logger, client: client}
	return dao
}

// PutMemberJobs ..
func (dao *jobDao) PutMemberJobIDs(nodeid string, jobIDs []string) (err error) {
	jobIDsBytes, err := dao.client.MarshalObject(jobIDs)
	
	if err != nil {
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathMemberJobs, nodeid, jobIDsBytes)
	
	return dao.client.BroadcastTxSync(msg)
}

// GetMemberJobs ..
func (dao *jobDao) GetMemberJobIDs(nodeid string) (jobIDs []string, err error) {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathMemberJobs, nodeid)
	
	jobIDs = []string{}
	err = dao.client.GetObject(msg, jobIDs)
	return jobIDs, err
}

// GetMemberJobs ..
func (dao *jobDao) GetMemberJobs(membID string) (jobs []Job, err error) {
	jobIDs, err := dao.GetMemberJobIDs(membID)
	if err != nil {
		dao.logger.Error("[ERROR] Cannot retrieve member jobs", err)
		return []Job{}, err
	}
	jobs = []Job{}
	for _, jobID := range jobIDs {
		job, err2 := dao.GetJob(jobID)
		if err2 == nil {
			jobs = append(jobs, job)
		}
	}
	
	return jobs, err
}

// GetAllMemberJobIDs : returns member-JobIDs Map
func (dao *jobDao) GetAllMemberJobIDs() (membJobMap map[string][]string, err error) {
	msg := types.NewViewMsgMany(common.SpaceDaemon, PathMemberJobs, "", "")
	
	membJobMap = make(map[string][]string)
	
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		jobIDs := []string{}
		err := dao.client.UnmarshalObject(value, &jobIDs)
		if err != nil {
			dao.logger.Error("[ERROR-JobDao] unmarshal member jobs ", err)
		}
		membid := string(key)
		membJobMap[membid] = jobIDs
		return true
	})
	
	return membJobMap, err
}

// PutJob ..
func (dao *jobDao) PutJob(job Job) (err error) {
	bytes, err := dao.client.MarshalObject(job)
	
	if err != nil {
		dao.logger.Error("PutJob marshal", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathJobs, job.ID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

// RemoveJob ..
func (dao *jobDao) RemoveJob(jobID string) (err error) {
	msg := types.NewTxMsg(types.TxDelete, common.SpaceDaemon, PathJobs, jobID, nil)
	err = dao.client.BroadcastTxSync(msg)
	return err
}

// GetJob ..
func (dao *jobDao) GetJob(jobID string) (job Job, err error) {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathJobs, jobID)
	job = Job{}
	err = dao.client.GetObject(msg, &job)
	return job, err
}

// ContainsJob ..
func (dao *jobDao) ContainsJob(jobID string) bool {
	msg := types.NewViewMsgHas(common.SpaceDaemon, PathJobs, jobID)
	ok, err := dao.client.Has(msg)
	
	if err != nil {
		dao.logger.Error("[ERROR-JobDao] ContainsJob ", err)
	}
	return ok
}

// GetAllJobIDs ..
func (dao *jobDao) GetAllJobIDs() (jobIDs []string, err error) {
	msg := types.NewViewMsgKeys(common.SpaceDaemon, PathJobs, "", "")
	jobIDs, err = dao.client.GetKeys(msg)
	return jobIDs, err
}

// GetAllJobs ..
func (dao *jobDao) GetAllJobs() (jobs map[string]Job, err error) {
	msg := types.NewViewMsgKeys(common.SpaceDaemon, PathJobs, "", "")
	
	jobs = make(map[string]Job)
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		jobid := string(key)
		job := Job{}
		dao.client.UnmarshalObject(value, job)
		jobs[jobid] = job
		return true
	})
	
	return jobs, err
}


