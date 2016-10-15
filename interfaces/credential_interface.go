package interfaces

type ICredential interface {
	GetName() string
	GetPluginName() string
	Configure(name string, params map[string]interface{}) error
	Initialize() error
	GetParams() map[string]interface{}
	GetParam(paramName string) interface{}
}
