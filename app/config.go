package app

type Config struct {
	Collectors  []ConfigCollector `json:"collectors"`
	Credentials []ConfigCredential `json:"credentials"`
	DataSources []ConfigDataSource `json:"datasources"`
	Server      ConfigServer
}

type ConfigCollector struct {
	Name   string `json:"name"`
	Plugin string `json:"plugin"`
	Params map[string]interface{} `json:"params"`
}

type ConfigDataSource struct {
	Name   string `json:"name"`
	Plugin string `json:"plugin"`
	Params map[string]interface{} `json:"params"`
}

type ConfigCredential struct {
	Name   string `json:"name"`
	Plugin string `json:"plugin"`
	Params map[string]interface{} `json:"params"`
}

type ConfigServer struct {
	Port string `json:"port"`
	Host string `json:"host"`
}