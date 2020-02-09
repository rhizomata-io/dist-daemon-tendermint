package worker

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathCheckpoint = "chkpnt"
	
	// PatternPathData : data/{jobID}/{topic}
	PatternPathJobTopicData = "data/%s/%s"
)

// workerDao kv store model for cluster
type workerDao struct {
	config common.DaemonConfig
	logger log.Logger
	client types.Client
}

var _ Repository = (*workerDao)(nil)

func NewRepository(config common.DaemonConfig, logger log.Logger, client types.Client) Repository {
	dao := &workerDao{config: config, logger: logger, client: client}
	return dao
}

// PutCheckpoint marshal checkpoint to json
func (dao *workerDao) PutCheckpoint(jobID string, checkpoint interface{}) error {
	jobIDsBytes, err := dao.client.MarshalJson(checkpoint)
	if err != nil {
		return err
	}
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathCheckpoint, jobID, jobIDsBytes)
	return dao.client.BroadcastTxSync(msg)
}

// GetCheckpoint  unmarshal checkpoint from json
func (dao *workerDao) GetCheckpoint(jobID string, checkpoint interface{}) error {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathCheckpoint, jobID)
	bytes, err := dao.client.Query(msg)
	if err != nil {
		return err
	}
	err = dao.client.UnmarshalJson(bytes, &checkpoint)
	return err
}

func (path DataPath) String() string {
	return makeDataPath(path.JobID, path.Topic)
}

func makeDataPath(jobID string, topic string) string {
	return fmt.Sprintf(PatternPathJobTopicData, jobID, topic)
}

func ParseDataPathBytes(path []byte) (dataPath DataPath, err error) {
	paths := bytes.Split(path, []byte("/"))
	if len(paths) != 2 {
		return dataPath, errors.New(fmt.Sprintf("Illegal DataPath format %s", path))
	}
	dataPath = DataPath{JobID: string(paths[0]), Topic: string(paths[1])}
	return dataPath, err
}

func ParseDataPath(path string) (dataPath DataPath, err error) {
	paths := strings.Split(path, "/")
	if len(paths) != 2 {
		return dataPath, errors.New(fmt.Sprintf("Illegal DataPath format %s", path))
	}
	dataPath = DataPath{JobID: paths[0], Topic: paths[1]}
	return dataPath, err
}

// PutData put data to {common.SpaceDaemonData}
func (dao *workerDao) PutData(jobID string, topic string, rowID string, data []byte) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemonData, fullPath, rowID, data)
	return dao.client.BroadcastTxAsync(msg)
}

// PutDataFullPath ..
// func (dao *workerDao) PutDataFullPath(fullPath string, data string) error {
// 	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemonData, fullPath, "", data)
// 	_, err := dao.kv.Put(fullPath, data)
// 	if err != nil {
// 		log.Println("[ERROR-WorkerDao] PutDataFullPath", err)
// 	}
// 	return err
// }

// PutObject data type must be registered to Codec
func (dao *workerDao) PutObject(jobID string, topic string, rowID string, data interface{}) error {
	dataBytes, err := dao.client.MarshalObject(data)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] PutObject marshal", err)
	}
	// fmt.Println("&&&&& PutObject :: key=", key, ", data=", data)
	return dao.PutData(jobID, topic, rowID, dataBytes)
}

// PutObjectFullPath ..
// func (dao *workerDao) PutObjectFullPath(fullPath string, data interface{}) error {
// 	_, err := dao.kv.PutObject(fullPath, data)
// 	if err != nil {
// 		log.Println("[ERROR-WorkerDao] PutObjectFullPath", err)
// 	}
// 	// fmt.Println("&&&&& PutObjectFullPath :: fullPath=", fullPath, ", data=", data)
// 	return err
// }

// GetData ..
func (dao *workerDao) GetData(jobID string, topic string, rowID string) (data []byte, err error) {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewViewMsgOne(common.SpaceDaemonData, fullPath, rowID)
	bytes, err := dao.client.Query(msg)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// GetObject data type must be registered to Codec
func (dao *workerDao) GetObject(jobID string, topic string, rowID string, data interface{}) error {
	bytes, err := dao.GetData(jobID, topic, rowID)
	if err != nil {
		return err
	}
	return dao.client.UnmarshalObject(bytes, data)
}

// DeleteData ..
func (dao *workerDao) DeleteData(jobID string, topic string, rowID string) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxDelete, common.SpaceDaemonData, fullPath, rowID, nil)
	err := dao.client.BroadcastTxSync(msg)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] DeleteData ", err)
	}
	
	return err
}

// DeleteDataFullPath ..
// func (dao *workerDao) DeleteDataFullPath(key string) error {
// 	_, err := dao.kv.DeleteOne(key)
// 	if err != nil {
// 		log.Println("[ERROR-WorkerDao] GetData ", err)
// 	}
// 	// fmt.Println("^^^^^^ DeleteDataFullPath :: key=", key)
// 	return err
// }

// GetDataWithTopic ..
func (dao *workerDao) GetDataWithTopic(jobID string, topic string, handler DataHandler) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewViewMsgMany(common.SpaceDaemonData, fullPath, "", "")
	
	err := dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		fmt.Println("key=", string(key), "value=", string(value))
		// handler(jobID,topic,rowID, value)
		return true
	})
	
	return err
}

// WatchDataWithTopic ..
// func (dao *workerDao) WatchDataWithTopic(jobid string, topic string,
// 	handler func(eventType kv.EventType, fullPath string, rowID string, value []byte)) *kv.Watcher {
// 	key := fmt.Sprintf(kvPatternDataTopic, dao.cluster, jobid, topic)
// 	watcher := dao.kv.WatchWithPrefix(key, handler)
// 	// fmt.Println("&&&&& WatchDataWithTopic :: key=", key)
//
// 	return watcher
// }
//
// // WatchData ..
// func (dao *workerDao) WatchData(jobid string, topic string, rowID string,
// 	handler func(eventType kv.EventType, fullPath string, rowID string, value []byte)) *kv.Watcher {
// 	key := fmt.Sprintf(kvPatternData, dao.cluster, jobid, topic, rowID)
// 	watcher := dao.kv.Watch(key, handler)
// 	// fmt.Println("&&&&& WatchData :: key=", key)
//
// 	return watcher
// }
