package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rhizomata-io/dist-daemon-tendermint/daemon"
)

// TMService ..
type TMAPI struct {
	daemon *daemon.Daemon
}

func NewTMAPI(daemon *daemon.Daemon) (api *TMAPI) {
	api = &TMAPI{
		daemon: daemon,
	}
	return api
}


func (api *TMAPI) RelativePath() string {
	return "tm"
}

func (api *TMAPI) SetHandlers(group *gin.RouterGroup) {
}




