package cluster

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	MemberChangedEventPath = types.EventPath("member-changed")
	LeaderChangedEventPath = types.EventPath("leader-changed")
)

type MemberChangedEvent struct {
	common.DaemonEvent
	IsLeader     bool
	AliveMembers []*Member
}

func (event *MemberChangedEvent) Path() types.EventPath { return MemberChangedEventPath }

type LeaderChangedEvent struct {
	common.DaemonEvent
	IsLeader bool
	Leader   *Member
}

func (event *LeaderChangedEvent) Path() types.EventPath { return LeaderChangedEventPath }
