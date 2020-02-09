package daemon

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	StartedEventPath = types.EventPath("daemon-started")
)

type StartedEvent struct {
	common.DaemonEvent
}

func (event StartedEvent) Path() types.EventPath { return StartedEventPath }

