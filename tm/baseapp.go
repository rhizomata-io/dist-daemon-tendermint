package tm

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/kv"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
	"strconv"
	"strings"
)

const (
	ValidatorSetChangePrefix string = "val:"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
	ProtocolVersion version.Protocol = 0x1
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	var state State
	state.db = db
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		panic(err)
	}
	if len(stateBytes) == 0 {
		return state
	}
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		panic(err)
	}
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

type BaseApplication struct {
	types.BaseApplication
	DB    *dbm.GoLevelDB
	state State
	
	// validator set
	ValUpdates []types.ValidatorUpdate
	
	valAddrToPubKeyMap map[string]types.PubKey
	
	logger log.Logger
}

var _ abcitypes.Application = (*BaseApplication)(nil)

func NewBaseApplication(dbDir string) (kvapp *BaseApplication){
	db, err := dbm.NewGoLevelDB("bridge_chain", dbDir)
	if err != nil {
		panic(err)
	}
	
	state := loadState(db)
	
	return &BaseApplication{DB:db,
		state:state,
		valAddrToPubKeyMap: make(map[string]types.PubKey),
		logger:             log.NewNopLogger(),
	}
}

func (app *BaseApplication) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *BaseApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	res := types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion.Uint64(),
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
	return res
}

func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ValidatorSetChangePrefix)
}
// format is "val:pubkey!power"
// pubkey is a base64-encoded 32-byte ed25519 key
func (app *BaseApplication) execValidatorTx(tx []byte) types.ResponseDeliverTx {
	tx = tx[len(ValidatorSetChangePrefix):]
	
	//get the pubkey and power
	pubKeyAndPower := strings.Split(string(tx), "!")
	if len(pubKeyAndPower) != 2 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected 'pubkey!power'. Got %v", pubKeyAndPower)}
	}
	pubkeyS, powerS := pubKeyAndPower[0], pubKeyAndPower[1]
	
	// decode the pubkey
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid base64", pubkeyS)}
	}
	
	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}
	
	// update
	return app.updateValidator(types.Ed25519ValidatorUpdate(pubkey, power))
}

func (app *BaseApplication) updateValidator(v types.ValidatorUpdate) types.ResponseDeliverTx {
	key := []byte("val:" + string(v.PubKey.Data))
	
	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], v.PubKey.Data)
	
	if v.Power == 0 {
		// remove validator
		hasKey, err := app.state.db.Has(key)
		if err != nil {
			panic(err)
		}
		if !hasKey {
			pubStr := base64.StdEncoding.EncodeToString(v.PubKey.Data)
			return types.ResponseDeliverTx{
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf("Cannot remove non-existent validator %s", pubStr)}
		}
		app.state.db.Delete(key)
		delete(app.valAddrToPubKeyMap, string(pubkey.Address()))
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Error encoding validator: %v", err)}
		}
		app.state.db.Set(key, value.Bytes())
		app.valAddrToPubKeyMap[string(pubkey.Address())] = v.PubKey
	}
	
	// we only update the changes array if we successfully updated the tree
	app.ValUpdates = append(app.ValUpdates, v)
	
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *BaseApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if isValidatorTx(req.Tx) {
		// update validators in the merkle tree
		// and in app.ValUpdates
		return app.execValidatorTx(req.Tx)
	}
	
	var key, value []byte
	parts := bytes.Split(req.Tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = req.Tx, req.Tx
	}
	
	fmt.Println("* DeliverTx : ", string(key), "=", string(value))
	
	app.state.db.Set(prefixKey(key), value)
	app.state.Size++
	
	events := []types.Event{
		{
			Type: "app",
			Attributes: []kv.Pair{
				//{Key: []byte("creator"), Value: []byte("Cosmoshi Netowoko")},
				{Key: []byte("key"), Value: key},
			},
		},
	}
	
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

func (app *BaseApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *BaseApplication) Commit() abcitypes.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

func (app *BaseApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery){
	fmt.Println("* Query : ", reqQuery)
	
	if reqQuery.Path == "/val"{
		key := []byte("val:" + string(reqQuery.Data))
		value, err := app.state.db.Get(key)
		if err != nil {
			panic(err)
		}
		
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		return
	}
	
	if reqQuery.Prove {
		value, err := app.state.db.Get(prefixKey(reqQuery.Data))
		if err != nil {
			panic(err)
		}
		if value == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		
		return
	}
	
	resQuery.Key = reqQuery.Data
	value, err := app.state.db.Get(prefixKey(reqQuery.Data))
	if err != nil {
		panic(err)
	}
	if value == nil {
		resQuery.Log = "does not exist"
	} else {
		resQuery.Log = "exists"
	}
	resQuery.Value = value
	
	return resQuery
}

func (app *BaseApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			app.logger.Error("Error updating validators", "r", r)
		}
	}
	
	fmt.Println("* InitChain : ", req)
	
	return types.ResponseInitChain{}
}

func (app *BaseApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	// reset valset changes
	app.ValUpdates = make([]types.ValidatorUpdate, 0)
	
	for _, ev := range req.ByzantineValidators {
		if ev.Type == tmtypes.ABCIEvidenceTypeDuplicateVote {
			// decrease voting power by 1
			if ev.TotalVotingPower == 0 {
				continue
			}
			app.updateValidator(types.ValidatorUpdate{
				PubKey: app.valAddrToPubKeyMap[string(ev.Validator.Address)],
				Power:  ev.TotalVotingPower - 1,
			})
		}
	}
	
	//fmt.Println("* BeginBlock : ", req)
	return types.ResponseBeginBlock{}
}

func (app *BaseApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	//fmt.Println("* EndBlock : ", req)
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}


//---------------------------------------------
// update validators

func (app *BaseApplication) Validators() (validators []types.ValidatorUpdate) {
	itr, err := app.state.db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; itr.Valid(); itr.Next() {
		if isValidatorTx(itr.Key()) {
			validator := new(types.ValidatorUpdate)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	return
}

func MakeValSetChangeTx(pubkey types.PubKey, power int64) []byte {
	pubStr := base64.StdEncoding.EncodeToString(pubkey.Data)
	return []byte(fmt.Sprintf("val:%s!%d", pubStr, power))
}
