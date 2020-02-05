package cluster

import (
	"time"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	PathMember    = "member"
	PathHeartbeat = "heartbeat"
	PathLeader = "leader"
)

// DAO kv store model for job
type DAO struct {
	config config.DaemonConfig
	logger log.Logger
	client types.Client
}

func (dao *DAO) PutMember(member *Member) (err error) {
	bytes, err := types.BasicCdc.MarshalBinaryBare(member)
	
	if err != nil {
		dao.logger.Error("PutMember marshal", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, config.SpaceDaemon, PathMember, member.NodeID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

func (dao *DAO) GetMember(nodeID string) (member Member, err error) {
	msg := types.NewViewMsgOne(config.SpaceDaemon, PathMember, nodeID)
	
	member = Member{}
	err = dao.client.GetObject(msg, &member)
	return member, err
}

// PutLeader set leader
func (dao *DAO) PutLeader(leader string) (err error) {
	msg := types.NewTxMsg(types.TxSet, config.SpaceDaemon, PathLeader, "", []byte(leader))
	return dao.client.BroadcastTxSync(msg)
}

// GetLeader get leader id
func (dao *DAO) GetLeader() (leader string, err error) {
	msg := types.NewViewMsgOne(config.SpaceDaemon, PathLeader, "")
	data, err := dao.client.Query(msg)
	return string(data), err
}

func (dao *DAO) GetAllMembers() (members []*Member, err error) {
	msg := types.NewViewMsgMany(config.SpaceDaemon, PathMember, "", "")
	
	members = []*Member{}
	
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		member := Member{}
		err = dao.client.UnmarshalObject(value, &member)
		if err != nil {
			dao.logger.Error("GetAllMembers unmarshal member : ", err)
		} else {
			members = append(members, &member)
		}
		
		return true
	})
	
	return members, err
}


func (dao *DAO) PutHeartbeat(nodeID string) (err error) {
	bytes, err := types.BasicCdc.MarshalBinaryBare(time.Now())
	
	if err != nil {
		dao.logger.Error("PutHeartbeat : Member : ", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, config.SpaceDaemon, PathHeartbeat, nodeID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

func (dao *DAO) GetHeartbeats(handler func(nodeid string, time time.Time)) (err error) {
	msg := types.NewViewMsgMany(config.SpaceDaemon, PathHeartbeat, "", "")
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		time := time.Time{}
		nodeid := string(key)
		err = dao.client.UnmarshalObject(value, &time)
		
		if err != nil {
			dao.logger.Error("GetHeartbeats unmarshal time", "key", string(key),
				"value", string(value), err)
		} else {
			handler(nodeid, time)
		}
		
		return true
	})
	
	return err
}
