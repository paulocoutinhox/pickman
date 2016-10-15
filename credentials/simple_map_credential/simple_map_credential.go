package credentials

import (
	"github.com/prsolucoes/pickman/managers"
	"reflect"
)

const (
	SIMPLE_MAP_DATA_CREDENTIAL_PLUGIN_NAME = "simple.map"
)

type SimpleMapDataCredential struct {
	Name   string
	Params map[string]interface{}
}

func init() {
	managers.CredentialsAvailable[SIMPLE_MAP_DATA_CREDENTIAL_PLUGIN_NAME] = reflect.TypeOf(SimpleMapDataCredential{});
}

func (This *SimpleMapDataCredential) GetName() string {
	return This.Name
}

func (This *SimpleMapDataCredential) GetPluginName() string {
	return SIMPLE_MAP_DATA_CREDENTIAL_PLUGIN_NAME
}

func (This *SimpleMapDataCredential) Configure(name string, params map[string]interface{}) error {
	This.Name = name
	This.Params = params
	return nil
}

func (This *SimpleMapDataCredential) Initialize() error {
	return nil
}

func (This *SimpleMapDataCredential) GetParams() map[string]interface{} {
	return This.Params
}

func (This *SimpleMapDataCredential) GetParam(paramName string) interface{} {
	param, ok := This.Params[paramName]

	if !ok {
		return ""
	}

	return param
}
