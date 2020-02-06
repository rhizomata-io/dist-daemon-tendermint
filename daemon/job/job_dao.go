package job

import (
	"encoding/json"
	"fmt"
	"github.com/rhizomata-io/dist-daemonize/kernel/kv"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathJobs    = "jobs"
	PathMemberJobs = "membjobs"
)

// DAO kv store model for job
type DAO struct {
	config config.DaemonConfig
	logger log.Logger
	client types.Client
}


// PutMemberJobs ..
func (dao *DAO) PutMemberJobs(nodeid string, jobIDs []string) (err error) {
	jobIDsBytes, err := dao.client.MarshalObject(jobIDs)
	
	if err != nil {
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, config.SpaceDaemon, PathMemberJobs, nodeid, jobIDsBytes)
	
	return dao.client.BroadcastTxSync(msg)
}

// GetMemberJobs ..
func (dao *DAO) GetMemberJobs(nodeid string) (jobIDs []string, err error) {
	msg := types.NewViewMsgOne(config.SpaceDaemon, PathMemberJobs, nodeid)
	
	jobIDs = []string{}
	err = dao.client.GetObject(msg, jobIDs)
	return jobIDs, err
}

// GetAllMemberJobIDs : returns member-JobIDs Map
func (dao *DAO) GetAllMemberJobIDs() (membJobMap map[string][]string, err error) {
	membJobMap = make(map[string][]string)
	dirPath := fmt.Sprintf(kvDirMemberJob, dao.cluster)
	err = dao.client .GetWithPrefix(dirPath,
		func(fullPath string, rowID string, value []byte) bool {
			jobIDs := []string{}
			err := json.Unmarshal(value, &jobIDs)
			if err != nil {
				dao.logger.Info("[ERROR-JobDao] unmarshal member jobs ", fullPath, err)
			}
			membid := rowID
			membJobMap[membid] = jobIDs
			return true
		})
	
	return membJobMap, err
}

// WatchMemberJobs ..
func (dao *DAO) WatchMemberJobs(memberID string, handler func(jobIDs []string)) (watcher *kv.Watcher) {
	dirPath := fmt.Sprintf(kvPatternMemberJob, dao.cluster, memberID)
	watcher = dao.client .Watch(dirPath,
		func(eventType kv.EventType, fullPath string, rowID string, value []byte) {
			jobIDs := []string{}
			err := json.Unmarshal(value, &jobIDs)
			if err != nil {
				dao.logger.Info("[ERROR-JobDao] unmarshal member jobs ", memberID, err)
			}
			handler(jobIDs)
		})
	return watcher
}

// GetJob ..
func (dao *DAO) GetJob(jobID string) (job Job, err error) {
	value, err := dao.client .GetOne(fmt.Sprintf(kvPatternJob, dao.cluster, jobID))
	return Job{ID: jobID, Data: value}, err
}

// ContainsJob ..
func (dao *DAO) ContainsJob(jobID string) bool {
	value, err := dao.client .GetOne(fmt.Sprintf(kvPatternJob, dao.cluster, jobID))
	return err == nil && value != nil
}

// PutJob ..
func (dao *DAO) PutJob(jobID string, value []byte) (err error) {
	_, err = dao.client .Put(fmt.Sprintf(kvPatternJob, dao.cluster, jobID), string(value))
	return err
}

// RemoveJob ..
func (dao *DAO) RemoveJob(jobID string) (err error) {
	_, err = dao.client .DeleteOne(fmt.Sprintf(kvPatternJob, dao.cluster, jobID))
	return err
}

// GetAllJobIDs ..
func (dao *DAO) GetAllJobIDs() (jobIDs []string, err error) {
	jobIDs = []string{}
	dirPath := fmt.Sprintf(kvPatternJobsDir, dao.cluster)
	err = dao.client .GetWithPrefix(dirPath,
		func(fullPath string, rowID string, value []byte) bool {
			jobid := rowID
			jobIDs = append(jobIDs, jobid)
			return true
		})
	
	return jobIDs, err
}

// GetAllJobs ..
func (dao *DAO) GetAllJobs() (jobs map[string]Job, err error) {
	jobs = make(map[string]Job)
	dirPath := fmt.Sprintf(kvPatternJobsDir, dao.cluster)
	err = dao.client .GetWithPrefix(dirPath,
		func(fullPath string, rowID string, value []byte) bool {
			jobid := rowID
			job := Job{ID: jobid, Data: value}
			jobs[jobid] = job
			return true
		})
	
	return jobs, err
}

// WatchJobs ..
func (dao *DAO) WatchJobs(handler func(jobid string, data []byte)) (watcher *kv.Watcher) {
	dirPath := fmt.Sprintf(kvPatternJobsDir, dao.cluster)
	watcher = dao.client .WatchWithPrefix(dirPath,
		func(eventType kv.EventType, fullPath string, rowID string, value []byte) {
			jobid := rowID
			handler(jobid, value)
		})
	return watcher
}

