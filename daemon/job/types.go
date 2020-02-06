package job

import (
	"encoding/json"
	
	"github.com/google/uuid"
)

const (
	sharp = byte('#')
	colon = byte(':')
)

// Job job data structure
type Job struct {
	ID   string
	Data []byte
}

//PIandBody : devides data into Processing Instartuction('#xxx:') and Body
type PIandBody struct {
	PI   string
	Body []byte
}

// NewJob create new job with uuid
func NewJob(data []byte) Job {
	uuid := uuid.New()
	return Job{ID: uuid.String(), Data: data}
}

// NewWithPI create new job with pi
func NewWithPI(pi string, data string) Job {
	data = "#" + pi + ":" + data
	uuid := uuid.New()
	return Job{ID: uuid.String(), Data: []byte(data)}
}

// NewWithPIAndID create new job with pi
func NewWithPIAndID(jobID string, pi string, data string) Job {
	data = "#" + pi + ":" + data
	return Job{ID: jobID, Data: []byte(data)}
}

// NewWithID create new job with pi
func NewWithID(jobID string, data []byte) Job {
	return Job{ID: jobID, Data: data}
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

// HasPI check whether data has PI
func (job *Job) HasPI() bool {
	if job.Data[0] == sharp {
		index := -1
		for i, b := range job.Data {
			index = i
			if b == colon {
				break
			}
		}
		
		if index > 1 {
			return true
		}
	}
	return false
}

// GetPIandBody get PIandBody if Data has PI
func (job *Job) GetPIandBody() (pib *PIandBody) {
	if job.Data[0] == sharp {
		index := -1
		for i, b := range job.Data {
			index = i
			if b == colon {
				break
			}
		}
		
		if index < 2 {
			return nil
		}
		
		pi := string(job.Data[1:index])
		body := job.Data[index+1:]
		return &PIandBody{PI: pi, Body: body}
	}
	return nil
}
