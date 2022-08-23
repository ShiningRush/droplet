package wrapper

import (
	"reflect"

	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/middleware"
)

type WrapOptBase struct {
	inputType            reflect.Type
	isReadFromBody       bool
	disableUnmarshalBody bool

	trafficLogOpt middleware.TrafficLogOpt

	orchestrator droplet.Orchestrator
}

type SetWrapOpt func(base *WrapOptBase)

func InputType(p reflect.Type) SetWrapOpt {
	return func(base *WrapOptBase) {
		if p.Kind() == reflect.Ptr {
			p = p.Elem()
		}
		if p.Kind() != reflect.Struct {
			panic("input type must be struct or struct ptr")
		}
		base.inputType = p
	}
}

func ReadFromBody() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.isReadFromBody = true
	}
}

func LogReqAndResp() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.trafficLogOpt.IsLogReqAndResp = true
	}
}

func LogFunc(logFunc func(log *middleware.TrafficLog)) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.trafficLogOpt.LogFunc = logFunc
	}
}

func Orchestrator(o droplet.Orchestrator) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.orchestrator = o
	}
}

func DisableUnmarshalBody() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.disableUnmarshalBody = true
	}
}
