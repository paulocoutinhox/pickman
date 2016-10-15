package managers

import (
	"github.com/prsolucoes/pickman/interfaces"
	"errors"
	"reflect"
)

var (
	Credentials []interfaces.ICredential
	CredentialsAvailable map[string]reflect.Type
)

func init() {
	Credentials = []interfaces.ICredential{}
	CredentialsAvailable = map[string]reflect.Type{}
}

func GetCredentialByName(name string) (interfaces.ICredential, error) {
	for _, credential := range Credentials {
		if credential.GetName() == name {
			return credential, nil
		}
	}

	return nil, errors.New("Credential not found")
}

func GetCredentialAvailableByPluginName(name string) (interfaces.ICredential, error) {
	pluginType, ok := CredentialsAvailable[name]

	if !ok {
		return nil, errors.New("Credential available not found")
	}

	pluginRef := reflect.New(pluginType);
	plugin := pluginRef.Interface().(interfaces.ICredential)

	return plugin, nil
}