package wrapper

import (
	"reflect"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/middleware"
)

type WrapOptBase struct {
	inputType            reflect.Type
	isReadFromBody       bool
	disableUnmarshalBody bool

	trafficLogOpt middleware.TrafficLogOpt

	orchestrator core.Orchestrator
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

func LogReq() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.trafficLogOpt.LogReq = true
	}
}

func LogResp() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.trafficLogOpt.LogResp = true
	}
}

func SetLogger(logger middleware.TrafficLogger) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.trafficLogOpt.Logger = logger
	}
}

func Orchestrator(o core.Orchestrator) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.orchestrator = o
	}
}

func DisableUnmarshalBody() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.disableUnmarshalBody = true
	}
}
