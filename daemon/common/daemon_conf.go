package common

import "github.com/rhizomata-io/dist-daemon-tendermint/types"

const (
	SpaceDaemon = "daemon"
)

type DaemonConfig struct {
	ChainID  string
	NodeID   string
	NodeName string
	// HeartbeatInterval Heartbeat Interval
	HeartbeatInterval uint
	
	// CheckHeartbeatInterval Heartbeat check Interval
	// CheckHeartbeatInterval uint
	
	// AliveThreasholdSecond Heartbeat time Threshold
	AliveThresholdSeconds uint
}
