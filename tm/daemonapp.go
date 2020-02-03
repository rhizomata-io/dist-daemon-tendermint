package tm

import (
	"fmt"
	
	"github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type DaemonApp struct {
	abcitypes.BaseApplication
}

func NewDaemonApplication(dbDir string) (dapp *DaemonApp) {
	return &DaemonApp{}
}

func (app *DaemonApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	msg, err := types.Decode(req.Tx)
	
	if err != nil {
		fmt.Println("[ERROR] DeliverTx", err, string(req.Tx))
	} else {
		var ptr interface{}
		switch msg.Path {
		case types.PathHeartbeat:
			ptr = &types.Member{}
		default:
			panic("Unknown message type")
		}
		
		err = types.BasicCdc.UnmarshalBinaryBare(msg.Data, ptr)
		
		if err != nil {
			fmt.Println("Unknown message data ", err)
		}
		// fmt.Println("Data::", msg.Data)
		fmt.Println("ptr::", ptr)
	}
	
	
	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *DaemonApp) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery){
	fmt.Println("* Query : reqQuery.Path=", reqQuery.Path)
	fmt.Println("          reqQuery.Data=", string(reqQuery.Data))
	return resQuery
}