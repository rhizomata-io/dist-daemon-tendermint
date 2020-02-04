package daemon

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	"time"
)

func (dm *Daemon)PutHeartbeat() (err error){
	member := types.NewMember(dm.config.Moniker, string(dm.tmNode.NodeInfo().ID()))
	member.Heartbeat = time.Now()
	
	bytes, err := types.BasicCdc.MarshalBinaryBare(member)
	
	if err != nil {
		dm.logger.Error("MarshalBinaryBare : Member : ", err)
		return err
	}
	
	msg := types.NewTxMsg(types.MethodPutHeartbeat, types.PathMember, member.NodeID, bytes)
	
	return dm.client.BroadcastTxSync(msg)
}

func (dm *Daemon)GetMember(nodeID string) (member *types.Member, err error){
	msg := types.NewViewMsgOne(types.PathMember, nodeID)
	resp, err := dm.client.Query(msg)
	if err == nil {
		member = &types.Member{}
		err = types.BasicCdc.UnmarshalBinaryBare(resp.Value, member)
		if err != nil {
			member = nil
			dm.logger.Error("GetMember : ", err)
		}
	}
	return member, err
}

func (dm *Daemon)GetAllMembers() (members []*types.Member, err error){
	msg := types.NewViewMsgMany(types.PathMember, "", "")
	
	resp, err := dm.client.Query(msg)
	if err != nil {
		dm.logger.Error("GetAllMembers : ", err)
		return nil, err
	}
	
	bytesArray := [][]byte{}
	err = types.BasicCdc.UnmarshalBinaryBare(resp.Value, &bytesArray)
	
	if err != nil {
		dm.logger.Error("GetAllMembers Unmarshal : ", err)
		return nil, err
	}
	
	members = []*types.Member{}
	for _, bytes := range bytesArray {
		member := &types.Member{}
		err = types.BasicCdc.UnmarshalBinaryBare(bytes, member)
		if err != nil {
			dm.logger.Error("GetAllMembers unmarshal member : ", err)
		}
		members = append(members, member)
	}
	return members, err
}

