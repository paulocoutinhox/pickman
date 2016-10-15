package interfaces

type ICollector interface {
	GetName() string
	GetPluginName() string
	Configure(name string, params map[string]interface{}) error
	Initialize() error
	Collect() error
	GetParams() map[string]interface{}
	GetParam(paramName string)interface{}
}
