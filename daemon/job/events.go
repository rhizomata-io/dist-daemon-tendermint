package job

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

const (
	JobsChangedEventPath = types.EventPath("jobs-changed")
	MemberJobsChangedEventPath = types.EventPath("msm_jobs-changed")
)

type MemberJobsChangedEvent struct {
	common.DaemonEvent
	NodeID string
	JobIDs []string
}

func (event *MemberJobsChangedEvent) Path() types.EventPath { return MemberJobsChangedEventPath }

type JobsChangedEvent struct {
	common.DaemonEvent
	JobIDs []string
}

func (event *JobsChangedEvent) Path() types.EventPath { return JobsChangedEventPath }

