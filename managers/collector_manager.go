package managers

import (
	"github.com/prsolucoes/pickman/interfaces"
	"errors"
	"reflect"
)

var (
	Collectors []interfaces.ICollector
	CollectorsAvailable map[string]reflect.Type
)

func init() {
	Collectors = []interfaces.ICollector{}
	CollectorsAvailable = map[string]reflect.Type{}
}

func GetCollectorByName(name string) (interfaces.ICollector, error) {
	for _, collector := range Collectors {
		if collector.GetName() == name {
			return collector, nil
		}
	}

	return nil, errors.New("Collector not found")
}

func GetCollectorAvailableByPluginName(name string) (interfaces.ICollector, error) {
	pluginType, ok := CollectorsAvailable[name]

	if !ok {
		return nil, errors.New("Collector available not found")
	}

	pluginRef := reflect.New(pluginType);
	plugin := pluginRef.Interface().(interfaces.ICollector)

	return plugin, nil
}