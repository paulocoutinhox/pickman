package managers

import (
	"github.com/prsolucoes/pickman/interfaces"
	"errors"
	"reflect"
)

var (
	DataSources []interfaces.IDataSource
	DataSourcesAvailable map[string]reflect.Type
)

func init() {
	DataSources = []interfaces.IDataSource{}
	DataSourcesAvailable = map[string]reflect.Type{}
}

func GetDataSourceByName(name string) (interfaces.IDataSource, error) {
	for _, datasource := range DataSources {
		if datasource.GetName() == name {
			return datasource, nil
		}
	}

	return nil, errors.New("DataSource not found")
}

func GetDataSourceAvailableByPluginName(name string) (interfaces.IDataSource, error) {
	pluginType, ok := DataSourcesAvailable[name]

	if !ok {
		return nil, errors.New("DataSource available not found")
	}

	pluginRef := reflect.New(pluginType);
	plugin := pluginRef.Interface().(interfaces.IDataSource)

	return plugin, nil
}