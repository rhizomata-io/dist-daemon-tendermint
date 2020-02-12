package api

import (
	"encoding/json"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon/job"
	
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// DaemonAPI ..
type DaemonAPI struct {
	daemon *daemon.Daemon
}

func NewDaemonAPI(daemon *daemon.Daemon) (api *DaemonAPI) {
	api = &DaemonAPI{
		daemon: daemon,
	}
	return api
}

func (api *DaemonAPI) RelativePath() string {
	return "daemon"
}

func (api *DaemonAPI) SetHandlers(group *gin.RouterGroup) {
	group.POST("job/add", api.addJob)
	group.DELETE("job", api.removeJob)
	group.GET("jobs", api.getJobs)
}



func (api *DaemonAPI) getJobs(context *gin.Context) {
	memberJobs, err := api.daemon.GetJobRepository().GetAllMemberJobIDs()
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	bytes, err := json.Marshal(memberJobs)
	
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write(bytes)
	context.Writer.Flush()
}


func (api *DaemonAPI) addJob(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	factory := context.Param("factory")
	jobID := context.Param("jobid")
	
	var j job.Job
	
	if len(jobID) > 0 {
		j = job.NewWithID(factory, jobID, data)
	} else {
		j = job.New(factory, data)
	}
	
	err = api.daemon.GetJobRepository().PutJob(j)
	
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write([]byte(j.ID))
	context.Writer.Flush()
}

func (api *DaemonAPI) removeJob(context *gin.Context) {
	jobID := context.Param("jobid")
	
	err := api.daemon.GetJobRepository().RemoveJob(string(jobID))
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	context.Writer.WriteString("ok")
	context.Writer.Flush()
}
