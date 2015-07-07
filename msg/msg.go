package msg

import (
	"encoding/json"
	"reflect"
)

var TypeMap map[string]reflect.Type

func init() {
	TypeMap = make(map[string]reflect.Type)

	t := func(obj interface{}) reflect.Type { return reflect.TypeOf(obj).Elem() }
	TypeMap["TunnelRequest"] = t((*TunnelRequest)(nil))
	TypeMap["TunnelReply"] = t((*TunnelReply)(nil))
	TypeMap["Ping"] = t((*Ping)(nil))
	TypeMap["Pong"] = t((*Pong)(nil))
}

type Message interface{}

type Envelope struct {
	Type    string
	Payload json.RawMessage
}

type TunnelRequest struct {
	Port int	
}

type TunnelReply struct {
	URI	string	
}


type Ping struct {
}

type Pong struct {
}
