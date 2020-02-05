package cluster

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"
	"fmt"
	"time"
)

type Manager struct {
	config              config.DaemonConfig
	logger              log.Logger
	cluster             *Cluster
	dao                 *DAO
	running             bool
}

func NewManager(config config.DaemonConfig, logger log.Logger, client types.Client) *Manager {
	cluster := newCluster(config.ChainID)
	dao := DAO{config:config, logger: logger, client: client}
	
	manager := &Manager{
		config:              config,
		logger:              logger,
		cluster:             cluster,
		dao:                 &dao,
	}
	
	localMemb := Member{
		NodeID:    config.NodeID,
		Name:      config.NodeName,
		heartbeat: time.Now(),
		leader:    false,
		alive:     true,
		local:     true,
	}
	
	cluster.putMember(&localMemb)
	cluster.localMember = &localMemb
	
	return manager
}

func (manager *Manager) Start() {
	manager.running = true
	
	err := manager.dao.PutMember(manager.cluster.localMember)
	
	if err != nil {
		manager.logger.Error("Cannot PutMember.", err)
		panic(err)
	}
	
	err = manager.dao.PutHeartbeat(manager.config.NodeID)
	
	if err != nil {
		manager.logger.Error("Cannot send Heartbeat.", err)
		panic(err)
	}
	
	go func() {
		for manager.running {
			time.Sleep(time.Duration(manager.config.HeartbeatInterval) * time.Second)
			err := manager.dao.PutHeartbeat(manager.config.NodeID)
			if err != nil {
				manager.logger.Error("Cannot send Heartbeat.", err)
			}
		}
	}()
	
	go func() {
		for manager.running {
			time.Sleep(time.Duration(manager.config.HeartbeatInterval + 1) * time.Second)
			changed := false
			err := manager.dao.GetHeartbeats(func (nodeid string, tm time.Time){
				c := manager.handleHeartbeat(nodeid, tm)
				if c {
					changed = true
				}
			})
			if err != nil {
				manager.logger.Error("[FATAL] Cannot check heartbeats.", err)
			}
			
			manager.checkLeader()
			
			if changed {
				manager.onMemberChanged()
			}
			
		}
	}()
	
	manager.logger.Info("[INFO-Cluster] Start Cluster Manager.")
}

// returns true if member state changed
func (manager *Manager)handleHeartbeat(nodeid string, tm time.Time) (changed bool) {
	member := manager.cluster.GetMember(nodeid)
	if member == nil {
		memb, err := manager.dao.GetMember(nodeid)
		if err != nil {
			manager.logger.Error(fmt.Sprintf("[FATAL] Cannot find member [%s]", nodeid), err)
			return false
		}
		member = &memb
		manager.cluster.putMember(member)
	}
	
	oldAlive := member.IsAlive()
	
	// fmt.Println("   -- ", nodeid , "1 oldAlive=", oldAlive , ", tm=" , tm, ", heartbeat=", member.Heartbeat())
	
	if member.IsLocal() {
		member.SetAlive(true)
	} else if member.Heartbeat().IsZero() || tm.Equal(member.Heartbeat()) {
		gap := time.Now().Sub(tm)
		// fmt.Println("   -- ", nodeid , "2 gap=", gap , ", tm=" , tm, ", heartbeat=", member.Heartbeat())
		if gap.Seconds() > float64(manager.config.AliveThresholdSeconds) {
			member.SetAlive(false)
		} else {
			member.SetAlive(true)
		}
	} else {
		member.SetAlive(true)
	}
	
	member.SetHeartbeat(tm)
	
	changed = oldAlive != member.IsAlive()
	
	// fmt.Println("   -- ", nodeid , "3 IsAlive=", member.IsAlive() , ", changed=" , changed)
	
	return changed
}


func (manager *Manager)checkLeader() {
	leaderID, err := manager.dao.GetLeader()
	if err != nil {
		manager.logger.Error("Get Leader ", err)
	}
	
	oldLeader := manager.cluster.leader
	
	fmt.Println("   -- 1 leaderID=", leaderID , " oldLeader=", oldLeader)
	
	if oldLeader != nil {
		if oldLeader.NodeID == leaderID {
			if oldLeader.IsAlive() {
				return
			}
		} else {
			oldLeader.SetLeader(false)
			manager.cluster.leader = nil
		}
	}
	
	var leader *Member
	
	if len(leaderID) > 0 {
		leader = manager.cluster.GetMember(leaderID)
		if leader == nil {
			leader = manager.electLeader()
		} else if !leader.IsAlive() {
			leader.SetLeader(false)
			leader = manager.electLeader()
		}
	} else {
		leader = manager.electLeader()
	}
	
	leader.SetLeader(true)
	
	manager.cluster.leader = leader
	
	manager.onLeaderChanged(leader)
}


func (manager *Manager) electLeader() *Member {
	members := manager.cluster.GetSortedMembers()
	
	fmt.Println("****** electLeader:: len(members) ", len(members))
	
	for _, id := range members {
		memb := manager.cluster.GetMember(id)
		fmt.Println("    ****** electLeader:: member ", id, memb)
		if memb.IsAlive() {
			manager.dao.PutLeader(id)
			return memb
		}
	}
	
	local := manager.cluster.Local()
	manager.dao.PutLeader(local.NodeID)
	return local
}


// IsLeaderNode : returns whether this kernel is leader.
func (manager *Manager) IsLeaderNode() bool {
	return manager.cluster.localMember.IsLeader()
}

func (manager *Manager) onLeaderChanged(leader *Member) {
	if manager.cluster.localMember == leader {
		manager.logger.Info("[INFO-Cluster] Leader changed. I'm the leader")
	} else {
		manager.logger.Info("[INFO-Cluster] Leader changed.", "leader",
			manager.cluster.leader.NodeID )
	}
	
	types.Publish(&LeaderChangedEvent{
		IsLeader: manager.cluster.localMember.IsLeader(),
		Leader:   leader,
	})
}

func (manager *Manager) onMemberChanged() {
	manager.logger.Info("[INFO-Cluster] Members changed.", "members",
		manager.cluster.GetAliveMemberIDs())
	
	types.Publish(&MemberChangedEvent{
		IsLeader: manager.cluster.localMember.IsLeader(),
		AliveMembers:   manager.cluster.GetAliveMembers(),
	})
}


// GetCluster get cluster
func (manager *Manager) GetCluster() *Cluster {
	return manager.cluster
}