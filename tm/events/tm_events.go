package events

import (
	"time"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	EndBlockEventType   = types.EventType("end-block")
	BeginBlockEventType = types.EventType("begin-block")
)

type BeginBlockEvent struct {
	Height int64
	Time time.Time
}

func (event BeginBlockEvent) Type() types.EventType { return BeginBlockEventType }

var _ types.Event = BeginBlockEvent{}

type EndBlockEvent struct {
	Height int64
	Size int
}

func (event EndBlockEvent) Type() types.EventType { return EndBlockEventType }

var _ types.Event = EndBlockEvent{}
