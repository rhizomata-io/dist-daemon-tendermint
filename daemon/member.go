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
	
	msg := types.NewTxMsg(types.TxSet, types.PathMember, member.NodeID, bytes)
	
	return dm.client.BroadcastTxSync(msg)
}

func (dm *Daemon)GetMember(nodeID string) (member types.Member, err error){
	msg := types.NewViewMsgOne(types.PathMember, nodeID)
	
	member = types.Member{}
	err = dm.client.GetObject(msg, &member)
	return member, err
}

func (dm *Daemon)GetAllMembers() (members []*types.Member, err error){
	msg := types.NewViewMsgMany(types.PathMember, "", "")
	
	members = []*types.Member{}
	
	err = dm.client.GetMany(msg, func(data []byte)bool {
		member := types.Member{}
		err = dm.client.UnmarshalObject(data, &member)
		if err != nil {
			dm.logger.Error("GetAllMembers unmarshal member : ", err)
		} else {
			members = append(members, &member)
		}
		
		return true
	})
	
	return members, err
}

