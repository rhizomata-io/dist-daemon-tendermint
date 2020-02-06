package job

import (
	"fmt"
	
	"github.com/rhizomata-io/dist-daemonize/kernel/kv"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathJobs       = "jobs"
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
	msg := types.NewViewMsgMany(config.SpaceDaemon, PathMemberJobs, "", "")
	
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

// // WatchMemberJobs ..
// func (dao *DAO) WatchMemberJobs(memberID string, handler func(jobIDs []string)) (watcher *kv.Watcher) {
// 	dirPath := fmt.Sprintf(kvPatternMemberJob, dao.cluster, memberID)
// 	watcher = dao.client .Watch(dirPath,
// 		func(eventType kv.EventType, fullPath string, rowID string, value []byte) {
// 			jobIDs := []string{}
// 			err := json.Unmarshal(value, &jobIDs)
// 			if err != nil {
// 				dao.logger.Info("[ERROR-JobDao] unmarshal member jobs ", memberID, err)
// 			}
// 			handler(jobIDs)
// 		})
// 	return watcher
// }

// PutJob ..
func (dao *DAO) PutJob(job Job) (err error) {
	bytes, err := dao.client.MarshalObject(job)
	
	if err != nil {
		dao.logger.Error("PutJob marshal", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, config.SpaceDaemon, PathJobs, job.ID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

// RemoveJob ..
func (dao *DAO) RemoveJob(jobID string) (err error) {
	msg := types.NewTxMsg(types.TxDelete, config.SpaceDaemon, PathJobs, jobID, nil)
	err = dao.client.BroadcastTxSync(msg)
	return err
}

// GetJob ..
func (dao *DAO) GetJob(jobID string) (job Job, err error) {
	msg := types.NewViewMsgOne(config.SpaceDaemon, PathJobs, jobID)
	job = Job{}
	err = dao.client.GetObject(msg, &job)
	return job, err
}

// ContainsJob ..
func (dao *DAO) ContainsJob(jobID string) bool {
	msg := types.NewViewMsgHas(config.SpaceDaemon, PathJobs, jobID)
	ok, err := dao.client.Has(msg)
	
	if err != nil {
		dao.logger.Error("[ERROR-JobDao] ContainsJob ", err)
	}
	return ok
}

// GetAllJobIDs ..
func (dao *DAO) GetAllJobIDs() (jobIDs []string, err error) {
	msg := types.NewViewMsgKeys(config.SpaceDaemon, PathJobs, "", "")
	jobIDs, err = dao.client.GetKeys(msg)
	return jobIDs, err
}

// GetAllJobs ..
func (dao *DAO) GetAllJobs() (jobs map[string]Job, err error) {
	msg := types.NewViewMsgKeys(config.SpaceDaemon, PathJobs, "", "")
	
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

// WatchJobs ..
// func (dao *DAO) WatchJobs(handler func(jobid string, data []byte)) (watcher *kv.Watcher) {
// 	dirPath := fmt.Sprintf(kvPatternJobsDir, dao.cluster)
// 	watcher = dao.client.WatchWithPrefix(dirPath,
// 		func(eventType kv.EventType, fullPath string, rowID string, value []byte) {
// 			jobid := rowID
// 			handler(jobid, value)
// 		})
// 	return watcher
// }
