package wrapper

import (
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/middleware"
	"reflect"
)

type WrapOptBase struct {
	InputType      reflect.Type
	IsReadFromBody bool

	TrafficLogOpt middleware.TrafficLogOpt

	Orchestrator droplet.Orchestrator
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
		base.InputType = p
	}
}

func ReadFromBody(p reflect.Type) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.IsReadFromBody = true
	}
}

func LogReqAndResp() SetWrapOpt {
	return func(base *WrapOptBase) {
		base.TrafficLogOpt.IsLogReqAndResp = true
	}
}

func LogFunc(logFunc func(log *middleware.TrafficLog)) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.TrafficLogOpt.LogFunc = logFunc
	}
}

func Orchestrator(o droplet.Orchestrator) SetWrapOpt {
	return func(base *WrapOptBase) {
		base.Orchestrator = o
	}
}
