package daemon

import (
	"time"
	
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type Daemon struct {
	tmNode *node.Node
	config *cfg.Config
	logger log.Logger
}

func NewDaemon(config *cfg.Config, logger log.Logger, tmNode *node.Node) (dm *Daemon) {
	dm = &Daemon{config:config, logger:logger, tmNode:tmNode}
	return dm
}

func (dm *Daemon) Start(){
	dm.logger.Info("Start Daemon..." , "node_id", dm.tmNode.NodeInfo().ID())
	
	
	go func(){
		time.Sleep(2*time.Second)
		for true {
			time.Sleep(2*time.Second)
	
			nodeID := string(dm.tmNode.NodeInfo().ID())
			member := types.NewMember(dm.config.Moniker, nodeID)
			member.Heartbeat = time.Now()
	
			bytes, err := types.BasicCdc.MarshalBinaryBare(member)
	
			if err != nil {
				dm.logger.Error("MarshalBinaryBare : Member : ", err)
			} else {
				msgBytes, err := types.Encode(types.PathHeartbeat, nodeID, bytes)
				if err != nil {
					dm.logger.Error("MarshalBinaryBare : Msg : ", err)
				} else {
					res, err := core.BroadcastTxSync(&rpctypes.Context{}, msgBytes)
					if err != nil {
						dm.logger.Error("BroadcastTxSync : Member : ", err)
						// fmt.Println("[ERROR] BroadcastTxSync ::", err)
					} else {
						fmt.Println("BroadcastTxSync Result::", res.Code)
					}
				}
			}
		}
	}()
	
	
	// go func(){
	// 	time.Sleep(2*time.Second)
	// 	for i:=0;i<100;i++{
	// 		time.Sleep(30*time.Millisecond)
	// 		stt , _ := core.Status(&rpctypes.Context{})
	// 		core.BroadcastTxSync(&rpctypes.Context{}, []byte(fmt.Sprintf("stest%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	// 		//core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("stest%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	//
	// 		if i%7 ==0{
	// 			time.Sleep(300*time.Millisecond)
	// 			core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("Commit%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	// 		}
	// 	}
	//
	// }()
	//
	// go func(){
	// 	time.Sleep(5*time.Second)
	// 	for i:=0;i<100;i++ {
	// 		time.Sleep(30*time.Millisecond)
	//
	// 		res, err := core.ABCIQuery(&rpctypes.Context{},"", []byte{}, 0, false)
	//
	// 		if err != nil {
	// 			dm.logger.Error("[ERROR] ABCIQuery ", err)
	// 		} else {
	// 			bytes, err := res.Response.MarshalJSON()
	// 			if err != nil {
	// 				dm.logger.Error("[ERROR] ABCIQuery Result MarshalJSON  ", err)
	// 			} else {
	// 				fmt.Println("ABCIQuery Result :: ", string(bytes))
	// 			}
	// 		}
	// 		//iterator, _ := provider.App.DB.Iterator([]byte("kvPairKey:stest11"),[]byte("kvPairKey:stest20"))
	// 		//for iterator.Valid() {
	// 		//	fmt.Println(" ^^ DB.Iterator: key=", string(iterator.Key()), ", value=", string(iterator.Value()))
	// 		//	iterator.Next()
	// 		//}
	// 		//iterator.Close()
	// 		//if i%7 ==0{
	// 		//	time.Sleep(300*time.Millisecond)
	// 		//}
	// 	}
	// }()
}