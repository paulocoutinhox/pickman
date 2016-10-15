package interfaces

type IDataSource interface {
	GetName() string
	GetPluginName() string
	Configure(name string, params map[string]interface{}) error
	Initialize() error
}
