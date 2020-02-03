package main

import (
	"github.com/rhizomata-io/dist-daemon-tendermint/tm"
	"os"
	"path/filepath"
	
	cmd "github.com/rhizomata-io/dist-daemon-tendermint/cmd/commands"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	DefaultBCDir = "chainroot"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.InitFilesCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		cmd.ShowNodeIDCmd,
		cmd.VersionCmd,
	)
	
	// Create & start node
	rootCmd.AddCommand(cmd.NewRunNodeCmd(tm.NewNode))
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	
	//go func(){
	//	time.Sleep(2*time.Second)
	//	for i:=0;i<100;i++{
	//		time.Sleep(20*time.Millisecond)
	//		stt , _ := core.Status(&rpctypes.Context{})
	//		core.BroadcastTxAsync(&rpctypes.Context{}, []byte(fmt.Sprintf("Async%d=%s%d",i,stt.NodeInfo.ID(), i)))
	//		//core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("test%d=%s%d",i,stt.NodeInfo.ID(), i)))
	//
	//		if i%5 ==0{
	//			core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("Commit%d=%s%d",i,stt.NodeInfo.ID(), i)))
	//		}
	//
	//		if i%6 ==0{
	//			time.Sleep(200*time.Millisecond)
	//		}
	//	}
	//}()
	//
	//
	//go func(){
	//	time.Sleep(2*time.Second)
	//	for i:=0;i<100;i++{
	//		time.Sleep(30*time.Millisecond)
	//		stt , _ := core.Status(&rpctypes.Context{})
	//		core.BroadcastTxSync(&rpctypes.Context{}, []byte(fmt.Sprintf("stest%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	//		//core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("stest%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	//
	//		if i%7 ==0{
	//			time.Sleep(300*time.Millisecond)
	//			core.BroadcastTxCommit(&rpctypes.Context{}, []byte(fmt.Sprintf("Commit%d=%s%ds",i,stt.NodeInfo.ID(), i)))
	//		}
	//	}
	//
	//}()
	//
	//go func(){
	//	time.Sleep(5*time.Second)
	//	for i:=0;i<100;i++ {
	//		time.Sleep(30*time.Millisecond)
	//		iterator, _ := provider.App.DB.Iterator([]byte("kvPairKey:stest11"),[]byte("kvPairKey:stest20"))
	//		for iterator.Valid() {
	//			fmt.Println(" ^^ DB.Iterator: key=", string(iterator.Key()), ", value=", string(iterator.Value()))
	//			iterator.Next()
	//		}
	//		iterator.Close()
	//		if i%7 ==0{
	//			time.Sleep(300*time.Millisecond)
	//		}
	//	}
	//}()
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
