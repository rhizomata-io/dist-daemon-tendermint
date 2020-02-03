package types

import (
	// "bytes"
	// "errors"
)

const(
	PathHeartbeat = "daemon/heartbeat"
	PathMember =  "daemon/member"
)


type Msg struct {
	Path string
	Key string
	Data []byte
}

func Encode(path string, key string, data []byte) ([]byte, error){
	msg := Msg{Path:path, Key:key, Data:data}
	return BasicCdc.MarshalBinaryBare(msg)
}

func Decode(msgBytes []byte) (*Msg, error) {
	msg := &Msg{}
	err := BasicCdc.UnmarshalBinaryBare(msgBytes, msg)
	return msg, err
}
