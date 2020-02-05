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
	*BaseApplication
	spaces map[string]*store.Registry
}

var _ abcitypes.Application = (*DaemonApp)(nil)

func NewDaemonApplication(config *cfg.Config, logger log.Logger, spaces []string) (dapp *DaemonApp) {
	baseapp := NewBaseApplication(config, logger)
	dapp = &DaemonApp{BaseApplication: baseapp, spaces: make(map[string]*store.Registry)}
	
	for _, name := range spaces {
		dapp.registerSpace(name)
	}
	
	return dapp
}

func (app *DaemonApp) registerSpace(name string) {
	db, err := dbm.NewGoLevelDB("daemon", app.config.DBDir())
	if err != nil {
		panic(err)
	}
	
	storeRegistry := store.NewRegistry(db)
	app.spaces[name] = storeRegistry
}

func (app *DaemonApp) getSpace(name string) *store.Registry {
	storeRegistry, ok := app.spaces[name]
	if !ok {
		panic(fmt.Sprintf("DB Space[%s] is not registered.", name))
	}
	return storeRegistry
}

func (app *DaemonApp) getStoreRegistryAny(space string, path string) *store.Store {
	storeRegistry := app.getSpace(space)
	return storeRegistry.GetOrMakeStore(path)
}

func (app *DaemonApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if app.isValidatorTx(req.Tx) {
		// update validators in the merkle tree
		// and in app.ValUpdates
		return app.execValidatorTx(req.Tx)
	}
	
	var cd uint32 = code.CodeTypeOK
	
	msg, err := types.DecodeTxMsg(req.Tx)
	
	if err != nil {
		app.logger.Error("DeliverTx", err)
		cd = code.CodeTypeEncodingError
	} else {
		store := app.getStoreRegistryAny(msg.Space, msg.Path)
		
		switch msg.Type {
		case types.TxSet:
			store.Set(msg.Key, msg.Data)
		case types.TxSetSync:
			store.SetSync(msg.Key, msg.Data)
		case types.TxDelete:
			store.Delete(msg.Key)
		case types.TxDeleteSync:
			store.DeleteSync(msg.Key)
		}
	}
	
	app.IncreaseTxSize()
	
	events := []abcitypes.Event{
		//{
		//	Type: "app",
		//	Attributes: []kv.Pair{
		//		{Key: []byte("key"), Value: key},
		//	},
		//},
	}
	
	return abcitypes.ResponseDeliverTx{Code: cd, Events: events}
}

func (app *DaemonApp) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	msg, err := types.DecodeViewMsg(reqQuery.Data)
	if err != nil {
		app.logger.Error("DecodeViewMsg ", err)
		resQuery.Log = err.Error()
		return resQuery
	}
	
	store, err := app.getSpace(msg.Space).GetStore(msg.Path)
	
	resQuery.Index = -1
	resQuery.Key = reqQuery.Data
	
	if err != nil {
		app.logger.Error("Unknown Store "+reqQuery.Path, err)
		resQuery.Log = err.Error()
		return resQuery
	}
	
	if msg.Type == types.GetOne {
		bytes, err := store.Get(msg.Start)
		if err != nil {
			app.logger.Error("Query GetOne "+string(msg.Start), err)
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
			app.logger.Error(fmt.Sprintf("Query GetMany %s - %s ", msg.Start, msg.End), err)
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
