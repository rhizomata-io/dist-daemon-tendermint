package daemon

import (
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/tendermint/tendermint/libs/log"
)

type DMClient struct {
	config *cfg.Config
	logger log.Logger
}


func NewClient( config *cfg.Config, logger log.Logger) types.Client {
	return &DMClient{ config, logger}
}

func (dm *DMClient) BroadcastTxSync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		dm.logger.Error("MarshalBinaryBare : Msg : ", err)
		return err
	}
	
	_, err = core.BroadcastTxSync(&rpctypes.Context{}, msgBytes)
	if err != nil {
		dm.logger.Error("BroadcastTxSync : ", err)
		return err
	}
	return err
}

func (dm *DMClient) BroadcastTxAsync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		dm.logger.Error("MarshalBinaryBare : Msg : ", err)
		return err
	}
	
	_, err = core.BroadcastTxAsync(&rpctypes.Context{}, msgBytes)
	if err != nil {
		dm.logger.Error("BroadcastTxAsync : ", err)
		return err
	}
	return err
}

func (dm *DMClient) Commit() (err error) {
	_, err = core.Commit(&rpctypes.Context{}, nil)
	if err != nil {
		dm.logger.Error("Commit : ", err)
	}
	return err
}


func (dm *DMClient) Query(msg *types.ViewMsg) (response abci.ResponseQuery, err error) {
	msgBytes, err := types.EncodeViewMsg(msg)
	if err != nil {
		dm.logger.Error("EncodeViewMsg : Msg : ", err)
		return response, err
	}
	
	res, err := core.ABCIQuery(&rpctypes.Context{}, msg.Path, bytes.HexBytes(msgBytes), 0, false)
	
	if err != nil {
		dm.logger.Error("Query : Msg : ", err)
	}
	
	return res.Response, err
}
