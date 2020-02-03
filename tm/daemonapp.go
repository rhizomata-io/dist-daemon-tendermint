package tm

import (
	"github.com/tendermint/tendermint/abci/types"
)


type DaemonApp struct {
	types.BaseApplication
	
}


func NewDaemonApplication(dbDir string) (dapp *DaemonApp){
	return &DaemonApp{ }
}
