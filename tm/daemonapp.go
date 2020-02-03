package tm

import (
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type DaemonApp struct {
	types.BaseApplication
}

func NewDaemonApplication(dbDir string) (dapp *DaemonApp) {
	return &DaemonApp{}
}

func (app *DaemonApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	fmt.Println("---- DeliverTx :: ", string(req.Tx))
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *DaemonApp) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery){
	fmt.Println("* Query : reqQuery.Path=", reqQuery.Path)
	fmt.Println("          reqQuery.Data=", string(reqQuery.Data))
	return resQuery
}