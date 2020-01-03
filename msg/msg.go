package msg

import (
	"reflect"

	"github.com/name5566/leaf/network/json"
)

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var Processor = json.NewProcessor()
var Processortype = reflect.TypeOf(Processor)

type RequestData struct {
	Route   string
	ReqID   uint64
	ReqData interface{}
}
type ResponseData struct {
	Code    int32
	Message string
	RspID   uint64
	RspName string
	RspData interface{}
}

func init() {
	// 这里我们注册了JSON 消息
	Processor.Register(&ResponseData{})
	Processor.Register(&RequestData{})
}
