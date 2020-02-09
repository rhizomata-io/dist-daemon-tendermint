package worker

import (
	"log"
	
	"github.com/rhizomata-io/dist-daemonize/kernel/kv"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathCheckpoint = "chkpnt"
)

// DAO kv store model for cluster
type DAO struct {
	config common.DaemonConfig
	logger log.Logger
	client types.Client
}

// PutCheckpoint ..
func (dao *DAO) PutCheckpoint(jobID string, checkpoint interface{}) error {
	jobIDsBytes, err := dao.client.MarshalJson(checkpoint)
	if err != nil {
		return err
	}
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathCheckpoint, jobID, jobIDsBytes)
	return dao.client.BroadcastTxSync(msg)
}

// GetCheckpoint ..
func (dao *DAO) GetCheckpoint(jobID string, checkpoint interface{}) error {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathCheckpoint, jobID)
	bytes, err := dao.client.Query(msg)
	if err != nil{
		return err
	}
	err = dao.client.UnmarshalJson(bytes, &checkpoint)
	return err
}

// PutData ..
func (dao *DAO) PutData(jobID string, topic string, rowID string, data string) error {
	fullPath := fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID)
	_, err := dao.kv.Put(fullPath, data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutData", err)
	}
	// fmt.Println("&&&&& PutData :: fullPath=", fullPath, ", data=", data)
	return err
}

// PutDataFullPath ..
func (dao *DAO) PutDataFullPath(fullPath string, data string) error {
	_, err := dao.kv.Put(fullPath, data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutDataFullPath", err)
	}
	return err
}

// PutObject ..
func (dao *DAO) PutObject(jobid string, topic string, rowID string, data interface{}) error {
	key := fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID)
	_, err := dao.kv.PutObject(key, data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutObject", err)
	}
	// fmt.Println("&&&&& PutObject :: key=", key, ", data=", data)
	return err
}

// PutObjectFullPath ..
func (dao *DAO) PutObjectFullPath(fullPath string, data interface{}) error {
	_, err := dao.kv.PutObject(fullPath, data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutObjectFullPath", err)
	}
	// fmt.Println("&&&&& PutObjectFullPath :: fullPath=", fullPath, ", data=", data)
	return err
}

// GetData ..
func (dao *DAO) GetData(jobid string, topic string, rowID string) (data []byte, err error) {
	data, err = dao.kv.GetOne(fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID))
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetData ", err)
	}
	return data, err
}

// GetObject ..
func (dao *DAO) GetObject(jobid string, topic string, rowID string, data interface{}) error {
	err := dao.kv.GetObject(fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID), data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetObject ", err)
	}
	return err
}

// DeleteData ..
func (dao *DAO) DeleteData(jobid string, topic string, rowID string) error {
	key := fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID)
	_, err := dao.kv.DeleteOne(key)
	if err != nil {
		log.Println("[ERROR-WorkerDao] DeleteData ", err)
	}
	// fmt.Println("^^^^^^ DeleteData :: key=", key)
	return err
}

// DeleteDataFullPath ..
func (dao *DAO) DeleteDataFullPath(key string) error {
	_, err := dao.kv.DeleteOne(key)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetData ", err)
	}
	// fmt.Println("^^^^^^ DeleteDataFullPath :: key=", key)
	return err
}

// GetDataWithTopic ..
func (dao *DAO) GetDataWithTopic(jobid string, topic string, handler kv.DataHandler) error {
	err := dao.kv.GetWithPrefix(fmt.Sprintf(kvPatternDataTopic, dao.cluster, jobid, topic), handler)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetDataWithJobID ", err)
	}
	return err
}

// WatchDataWithTopic ..
func (dao *DAO) WatchDataWithTopic(jobid string, topic string,
	handler func(eventType kv.EventType, fullPath string, rowID string, value []byte)) *kv.Watcher {
	key := fmt.Sprintf(kvPatternDataTopic, dao.cluster, jobid, topic)
	watcher := dao.kv.WatchWithPrefix(key, handler)
	// fmt.Println("&&&&& WatchDataWithTopic :: key=", key)

	return watcher
}

// WatchData ..
func (dao *DAO) WatchData(jobid string, topic string, rowID string,
	handler func(eventType kv.EventType, fullPath string, rowID string, value []byte)) *kv.Watcher {
	key := fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID)
	watcher := dao.kv.Watch(key, handler)
	// fmt.Println("&&&&& WatchData :: key=", key)

	return watcher
}
