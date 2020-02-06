package job

import (
	"encoding/json"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	
	"github.com/google/uuid"
)


func init() {
	types.BasicCdc.RegisterConcrete(Job{}, "daemon/job", nil)
}


// Job job data structure
type Job struct {
	WorkerFactory string
	ID            string
	Data          []byte
}

// NewJob create new job with uuid
func New(factory string, data []byte) Job {
	uuid := uuid.New()
	return Job{WorkerFactory: factory, ID: uuid.String(), Data: data}
}

// NewWithID create new job with pi
func NewWithID(factory string,  jobID string, data []byte) Job {
	return Job{WorkerFactory: factory, ID: jobID, Data: data}
}

// GetAsString Get data as string
func (job *Job) GetAsString() string {
	return string(job.Data)
}

// GetAsObject Get data as interface
func (job *Job) GetAsObject(obj interface{}) error {
	err := json.Unmarshal(job.Data, &obj)
	return err
}
