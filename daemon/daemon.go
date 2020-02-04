package daemon

import (
	"time"
	
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/tm/client"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type Daemon struct {
	tmNode *node.Node
	config *cfg.Config
	logger log.Logger
	client types.Client
}

func NewDaemon(config *cfg.Config, logger log.Logger, tmNode *node.Node) (dm *Daemon) {
	dm = &Daemon{config:config, logger:logger, tmNode:tmNode}
	dm.client = client.NewClient(config, logger, types.BasicCdc)
	return dm
}

func (dm *Daemon) Start(){
	dm.logger.Info("Start Daemon..." , "node_id", dm.tmNode.NodeInfo().ID())
	
	
	go func(){
		for true {
			time.Sleep(1*time.Second)
			dm.PutHeartbeat()
			
			member, err := dm.GetMember(string(dm.tmNode.NodeInfo().ID()))
			if err == nil {
				fmt.Println("----- After Heartbeat::", member)
			}
		}
	}()
	
	go func(){
		for i:=0;i<100;i++ {
			time.Sleep(2*time.Second)
			
			members, err := dm.GetAllMembers()
			if err != nil {
				fmt.Println("----- ERROR - GetAllMembers::", err)
			}
			
			fmt.Println("----- GetAllMembers::", members)
		}
	}()
}