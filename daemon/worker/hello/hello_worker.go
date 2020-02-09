package hello

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/worker"
)

const (FactoryName = "hello-worker")
type Factory struct {
}

var _ worker.Factory = (*Factory)(nil)

func (fac *Factory) Name() string { return FactoryName }

func (fac *Factory) NewWorker(helper *worker.Helper) (worker.Worker, error) {
	worker := &Worker{id: helper.ID(), helper: helper}
	
	jobInfo := JobInfo{}
	err := json.Unmarshal(worker.helper.Job().Data, &jobInfo)
	
	if err != nil {
		return nil, errors.New("cannot create hello worker " + err.Error())
	}
	
	interval, err := time.ParseDuration(worker.jobInfo.Interval)
	
	if err != nil {
		return nil, errors.New("cannot create hello worker " + err.Error())
	}
	
	jobInfo.interval = interval
	worker.jobInfo = jobInfo
	
	return worker, nil
}

type JobInfo struct {
	Interval string `json:"interval"`
	Greet    string `json:"greet"`
	interval time.Duration
}

type Worker struct {
	id      string
	started bool
	helper  *worker.Helper
	jobInfo JobInfo
}

var _ worker.Worker = (*Worker)(nil)

func (worker *Worker) ID() string      { return worker.id }
func (worker *Worker) IsStarted() bool { return worker.started }

func (worker *Worker) Start() error {
	fmt.Println(" - Worker Started ", worker.id)
	worker.started = true
	
	
	for worker.started {
		fmt.Printf( " - Worker[%s] worked. %s \n", worker.id, worker.jobInfo.Greet )
		time.Sleep(worker.jobInfo.interval)
	}
	
	fmt.Println(" - Worker Stopped ", worker.id)
	return nil
}

func (worker *Worker) Stop() error {
	worker.started = false
	return nil
}
