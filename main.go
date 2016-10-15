package main

import (
	"flag"
	"io/ioutil"
	"github.com/prsolucoes/pickman/app"
	"encoding/json"
	"github.com/prsolucoes/pickman/managers"
	"github.com/prsolucoes/pickman/logger"
	_ "github.com/prsolucoes/pickman/loaders"
)

var (
	Config *app.Config
	CollectorNameToExecute string
)

func main() {
	LoadConfigurations()
	LoadCredentialPlugins()
	LoadDataSourcePlugins()
	LoadCollectorPlugins()
	Execute()
}

func LoadConfigurations() {
	var configFileName = ""
	var configCollectorName = ""

	flag.StringVar(&configFileName, "c", "", "Example: -c config.json")
	flag.StringVar(&configCollectorName, "n", "", "Example: -n MyCollectorName")
	flag.Parse()

	if configFileName == "" {
		logger.F("You need set configuration file in JSON format. Type -h for help.")
	}

	if configCollectorName == "" {
		logger.F("You need set collector name to execute. Type -h for help.")
	}

	file, err := ioutil.ReadFile(configFileName)

	if err != nil {
		logger.F("Read file error (%v)", err)
	}

	// parse configuration file
	err = json.Unmarshal(file, &Config)

	if err != nil {
		logger.F("Parse file error (%v)", err)
	}

	// collector name to execute after load
	CollectorNameToExecute = configCollectorName
}

func Execute() {
	collector, err := managers.GetCollectorByName(CollectorNameToExecute)

	if err != nil {
		logger.F("Failed to get collector (%v - %v)", CollectorNameToExecute, err)
	}

	err = collector.Collect()

	if err != nil {
		logger.F("Failed to collect data (%v - %v)", CollectorNameToExecute, err)
	}

	logger.I("Data was collected with success")
}

func LoadCollectorPlugins() {
	if len(Config.Collectors) > 0 {
		for _, config := range Config.Collectors {
			collector, err := managers.GetCollectorAvailableByPluginName(config.Plugin)

			if err != nil {
				logger.F("Failed to load collector plugin (%v - %v)", config, err)
			}

			err = collector.Configure(config.Name, config.Params)

			if err != nil {
				logger.F("Failed to configure collector plugin (%v - %v)", config.Name, err)
			}

			logger.I("Collector plugin was configured (%v)", config.Name)

			// initialize plugin
			err = collector.Initialize()

			if err != nil {
				logger.F("Failed to initialize collector plugin (%v - %v)", config.Name, err)
			}

			logger.I("Collector plugin was initialized (%v)", config.Name)

			// add to collection
			managers.Collectors = append(managers.Collectors, collector)
		}
	} else {
		logger.I("You didn't configure any collector plugin in your configuration file")
	}
}

func LoadDataSourcePlugins() {
	if len(Config.DataSources) > 0 {
		for _, config := range Config.DataSources {
			datasource, err := managers.GetDataSourceAvailableByPluginName(config.Plugin)

			if err != nil {
				logger.F("Failed to load datasource plugin (%v - %v)", config, err)
			}

			err = datasource.Configure(config.Name, config.Params)

			if err != nil {
				logger.F("Failed to configure datasource plugin (%v - %v)", config.Name, err)
			}

			logger.I("DataSource plugin was configured (%v)", config.Name)

			// initialize plugin
			err = datasource.Initialize()

			if err != nil {
				logger.F("Failed to initialize datasource plugin (%v - %v)", config.Name, err)
			}

			logger.I("DataSource plugin was initialized (%v)", config.Name)

			// add to collection
			managers.DataSources = append(managers.DataSources, datasource)
		}
	} else {
		logger.I("You didn't configure any datasource plugin in your configuration file")
	}
}

func LoadCredentialPlugins() {
	if len(Config.Credentials) > 0 {
		for _, config := range Config.Credentials {
			credential, err := managers.GetCredentialAvailableByPluginName(config.Plugin)

			if err != nil {
				logger.F("Failed to load credential plugin (%v - %v)", config, err)
			}

			err = credential.Configure(config.Name, config.Params)

			if err != nil {
				logger.F("Failed to configure credential plugin (%v - %v)", config.Name, err)
			}

			logger.I("Credential plugin was configured (%v)", config.Name)

			// initialize plugin
			err = credential.Initialize()

			if err != nil {
				logger.F("Failed to initialize credential plugin (%v - %v)", config.Name, err)
			}

			logger.I("Credential plugin was initialized (%v)", config.Name)

			// add to collection
			managers.Credentials = append(managers.Credentials, credential)
		}
	} else {
		logger.I("You didn't configure any credential plugin in your configuration file")
	}
}