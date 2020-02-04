package tm

import (
	"fmt"
	"github.com/rhizomata-io/dist-daemon-tendermint/tm/store"
	cfg "github.com/tendermint/tendermint/config"
	dbm "github.com/tendermint/tm-db"
	
	"github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type DaemonApp struct {
	abcitypes.BaseApplication
	config *cfg.Config
	logger log.Logger
	storeRegistry *store.Registry
}

func NewDaemonApplication(config *cfg.Config, logger log.Logger, dbDir string) (dapp *DaemonApp) {
	db, err := dbm.NewGoLevelDB("daemon", dbDir)
	if err != nil {
		panic(err)
	}

	storeRegistry := store.NewRegistry(db)
	return &DaemonApp{config:config, logger:logger, storeRegistry:storeRegistry}
}

func (app *DaemonApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	var cd uint32 = code.CodeTypeOK
	
	msg, err := types.DecodeTxMsg(req.Tx)
	
	if err != nil {
		fmt.Println("[ERROR] DeliverTx", err, string(req.Tx))
		cd = code.CodeTypeEncodingError
	} else {
		store := app.storeRegistry.GetOrMakeStore(msg.Path)
		//fmt.Println("---- DeliverTx Store Path:", store.GetPath())
		store.Set(msg.Key, msg.Data)
	}
	
	return abcitypes.ResponseDeliverTx{Code: cd}
}

func (app *DaemonApp) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery){
	store, err := app.storeRegistry.GetStore(reqQuery.Path)
	
	resQuery.Index = -1
	resQuery.Key = reqQuery.Data
	
	if err != nil {
		app.logger.Error("Unknown Store " + reqQuery.Path, err)
		resQuery.Log = err.Error()
		return resQuery
	}
	
	msg, err := types.DecodeViewMsg(reqQuery.Data)
	if err != nil {
		app.logger.Error("DecodeViewMsg ", err)
		resQuery.Log = err.Error()
		return resQuery
	}
	
	if msg.Type == types.GetOne {
		bytes, err := store.Get(msg.Start)
		if err != nil {
			app.logger.Error("Query GetOne " + string(msg.Start), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if bytes == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		
		resQuery.Value = bytes
	} else if msg.Type == types.GetMany {
		bytes, err := store.GetMany(msg.Start, msg.End)
		
		if err != nil {
			app.logger.Error( fmt.Sprintf("Query GetMany %s - %s ",msg.Start, msg.End), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if bytes == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		
		resQuery.Value = bytes
	}
	
	return resQuery
}