package cluster

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	MemberChangedEventType = types.EventType("member-changed")
	LeaderChangedEventType = types.EventType("leader-changed")
)

type MemberChangedEvent struct {
	IsLeader     bool
	AliveMembers []*Member
}

func (event *MemberChangedEvent) Type() types.EventType { return MemberChangedEventType }

type LeaderChangedEvent struct {
	IsLeader bool
	Leader   *Member
}

func (event *LeaderChangedEvent) Type() types.EventType { return LeaderChangedEventType }
