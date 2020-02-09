package worker

import (
	"errors"
	"fmt"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/common"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

// Manager manager for jobs
type Manager struct {
	config  common.DaemonConfig
	dao     Repository
	logger  log.Logger
	facReg  *factoryRegistry
	workers map[string]Worker
}

// NewManager ..
func NewManager(config common.DaemonConfig, logger log.Logger, client types.Client) *Manager {
	dao := NewRepository(config, logger, client)
	manager := Manager{config: config, dao: dao, logger: logger}
	manager.facReg = NewFactoryRegistry()
	manager.workers = make(map[string]Worker)
	return &manager
}

func (manager *Manager) GetRepository() Repository {
	return manager.dao
}

func (manager *Manager) Start() {
}

// ContainsWorker if worker id is registered.
func (manager *Manager) ContainsWorker(id string) bool {
	return manager.workers[id] != nil
}

// GetWorker get worker for id
func (manager *Manager) GetWorker(id string) Worker {
	return manager.workers[id]
}

func (manager *Manager) newHelper(job job.Job) (helper *Helper) {
	helper = NewHelper(manager.config, manager.logger, job, manager.dao)
	return helper
}

// registerWorker ..
func (manager *Manager) registerWorker(job job.Job) error {
	if manager.workers[job.ID] != nil {
		return errors.New(fmt.Sprintf("worker[%s] is already registered. "+
			"If you want register new one, DeregisterWorker first", job.ID))
	}
	
	worker, err := manager.newWorker(job)
	if err != nil {
		return err
	}
	
	manager.workers[job.ID] = worker
	err = worker.Start()
	return err
}

func (manager *Manager) newWorker(job job.Job) (Worker, error) {
	helper := manager.newHelper(job)
	
	fac, err := manager.facReg.GetFactory(job.FactoryName)
	
	if err != nil {
		manager.logger.Error(fmt.Sprintf("cannot find worker factory '%s'", job.FactoryName), err)
		return nil, err
	}
	
	worker, err := fac.NewWorker(helper)
	
	if err != nil {
		manager.logger.Error("cannot create worker ", err)
		return nil, err
	}
	return worker, err
}

// deregisterWorker ..
func (manager *Manager) deregisterWorker(jobID string) error {
	worker := manager.workers[jobID]
	if worker == nil {
		return errors.New("Worker[" + jobID + "] is not registered.")
	}
	
	err := worker.Stop()
	
	if err == nil {
		delete(manager.workers, jobID)
	}
	
	return err
}

// SetJobs ...
func (manager *Manager) SetJobs(jobs []job.Job) {
	manager.logger.Info("[WorkerManager] Set Jobs:", len(jobs))
	
	tempWorkers := make(map[string]Worker)
	newWorkers := make(map[string]Worker)
	
	for id, worker := range manager.workers {
		tempWorkers[id] = worker
	}
	
	for _, job := range jobs {
		worker := tempWorkers[job.ID]
		if worker != nil {
			delete(tempWorkers, job.ID)
		} else {
			worker2, err := manager.newWorker(job)
			if err != nil {
				manager.logger.Error("[ERROR-WorkerMan] Cannot create worker ", err)
				continue
			} else {
				worker = worker2
				manager.logger.Info("[WARN-WorkerMan] New Worker..", job.ID)
			}
		}
		
		newWorkers[job.ID] = worker
	}
	// 제거된 worker 종료하기
	for id, worker := range tempWorkers {
		worker.Stop()
		manager.logger.Info("[WARN-WorkerMan] Dispose Worker .....", id)
	}
	
	manager.workers = newWorkers
	
	for id, worker := range manager.workers {
		if !worker.IsStarted() {
			go func(id string, worker Worker) {
				manager.logger.Info("[WARN-WorkerMan] New Worker Starting .....", id)
				worker.Start()
			}(id, worker)
		} else {
			manager.logger.Info("[WARN-WorkerMan] Remained Worker .....", id)
		}
	}
}