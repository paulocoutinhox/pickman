package collectors

import (
	"github.com/prsolucoes/pickman/managers"
	"github.com/prsolucoes/pickman/interfaces"
	"reflect"
	"errors"
	"fmt"
	"google.golang.org/api/analytics/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"net/http"
	"context"
	"io/ioutil"
	"github.com/prsolucoes/pickman/logger"
	"os/user"
	"path/filepath"
	"net/url"
	"encoding/json"
	"os"
	"strings"
	"strconv"
	"time"
)

const (
	GOOGLE_ANALYTICS_COLLECTOR_PLUGIN_NAME = "google.analytics"

	GOOGLE_ANALYTICS_COLLECTOR_AUTH_MODE_JWT = "jwt"
	GOOGLE_ANALYTICS_COLLECTOR_AUTH_MODE_OAUTH2 = "oauth2"

	GOOGLE_ANALYTICS_COLLECTOR_METRIC_TYPE_GA = "ga"
	GOOGLE_ANALYTICS_COLLECTOR_METRIC_TYPE_RT = "rt"
)

type GoogleAnalyticsCollector struct {
	Name       string
	Credential interfaces.ICredential
	Service    *analytics.Service
	Params     map[string]interface{}
}

func init() {
	managers.CollectorsAvailable[GOOGLE_ANALYTICS_COLLECTOR_PLUGIN_NAME] = reflect.TypeOf(GoogleAnalyticsCollector{});
}

func (This *GoogleAnalyticsCollector) GetName() string {
	return This.Name
}

func (This *GoogleAnalyticsCollector) GetPluginName() string {
	return GOOGLE_ANALYTICS_COLLECTOR_PLUGIN_NAME
}

func (This *GoogleAnalyticsCollector) Configure(name string, params map[string]interface{}) error {
	This.Name = name
	This.Params = params

	// check params because is required
	if params == nil {
		return errors.New("Params is empty")
	}

	// credential name
	credentialName, ok := params["credential"].(string)

	if !ok {
		return errors.New("Credential name is empty")
	}

	// credential
	credential, err := managers.GetCredentialByName(credentialName)

	if err != nil {
		return errors.New(fmt.Sprintf("Credential '%v' was not found", credentialName))
	}

	This.Credential = credential

	return nil
}

func (This *GoogleAnalyticsCollector) Collect() error {
	// validation
	viewId := This.GetParam("view.id").(string)

	if viewId == "" {
		return errors.New("Param view id is invalid")
	}

	metric := This.GetParam("metric").(string)

	if len(metric) < 3 {
		return errors.New("Param metric name is invalid")
	}

	metricData := strings.Split(metric, ":")

	if len(metricData) < 2 {
		return errors.New("Param metric name is invalid")
	}

	metricType := metricData[0]

	// execute by type
	if metricType == GOOGLE_ANALYTICS_COLLECTOR_METRIC_TYPE_GA {
		metricFrom := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		metricTo := time.Now().Format("2006-01-02")
		metricDimensions := "ga:operatingSystem"

		dataClient := This.Service.Data.Ga.Get(viewId, metricFrom, metricTo, metric)
		dataClient.Dimensions(metricDimensions)
		data, err := dataClient.Do()

		result := make(map[string]int64)

		if err == nil {
			for x := 0; x < len(data.Rows); x++ {
				item := data.Rows[x]
				source := item[0]
				total, err := strconv.ParseInt(item[1], 10, 64)

				if err != nil {
					return err
				}

				result[source] = total
			}
		} else {
			return err
		}

		logger.D("Data: %v", result)
	} else if metricType == GOOGLE_ANALYTICS_COLLECTOR_METRIC_TYPE_RT {
		dataClient := This.Service.Data.Realtime.Get(viewId, metric)
		data, err := dataClient.Do()

		if err != nil {
			return err
		}

		if len(data.Rows) > 0 {
			result, err := strconv.ParseInt(data.Rows[0][0], 10, 64)

			if err != nil {
				return err
			}

			logger.D("Data: %v", result)
		} else {
			return errors.New("Metric data not found")
		}
	} else {
		errors.New("Metric type is invalid")
	}

	return nil
}

func (This *GoogleAnalyticsCollector) Initialize() error {
	var client *http.Client
	ctx := context.Background()

	authMode := This.Credential.GetParam("auth.mode").(string)
	authFile := This.Credential.GetParam("auth.file").(string)

	if authMode == GOOGLE_ANALYTICS_COLLECTOR_AUTH_MODE_JWT {
		b, err := ioutil.ReadFile(authFile)

		if err != nil {
			return err
		}

		config, err := google.JWTConfigFromJSON(b, analytics.AnalyticsScope, analytics.AnalyticsReadonlyScope)

		if err != nil {
			return err
		}

		client = This.getClientJWT(ctx, config)
	} else if authMode == GOOGLE_ANALYTICS_COLLECTOR_AUTH_MODE_OAUTH2 {
		b, err := ioutil.ReadFile(authFile)

		if err != nil {
			return err
		}

		config, err := google.ConfigFromJSON(b, analytics.AnalyticsScope, analytics.AnalyticsReadonlyScope)

		if err != nil {
			return err
		}

		client, err = This.getClient(ctx, config)

		if err != nil {
			return err
		}
	} else {
		return errors.New("You need set the param: auth.mode")
	}

	s, err := analytics.New(client)

	if err != nil {
		return err
	}

	This.Service = s

	return nil
}

func (This *GoogleAnalyticsCollector) GetParams() map[string]interface{} {
	return This.Params;
}

func (This *GoogleAnalyticsCollector) GetParam(paramName string) interface{} {
	param, ok := This.Params[paramName]

	if !ok {
		return ""
	}

	return param
}

func (This *GoogleAnalyticsCollector) getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	// getClient uses a Context and Config to retrieve a Token
	// then generate a Client. It returns the generated Client.

	cacheFile, err := This.tokenCacheFile()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get path to cached credential file - %v", err))
	}

	token, err := This.tokenFromFile(cacheFile)

	if err != nil {
		token, err = This.getTokenFromWeb(config)

		if err != nil {
			return nil, err
		}

		This.saveToken(cacheFile, token)
	}

	return config.Client(ctx, token), nil
}

func (This *GoogleAnalyticsCollector) getClientJWT(ctx context.Context, config *jwt.Config) *http.Client {
	return config.Client(ctx)
}

func (This *GoogleAnalyticsCollector) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// getTokenFromWeb uses Config to request a Token.
	// It returns the retrieved Token.

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	logger.I("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var code string

	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	token, err := config.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (This *GoogleAnalyticsCollector) tokenCacheFile() (string, error) {
	// tokenCacheFile generates credential file path/filename.
	// It returns the generated credential path/filename.

	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".pickman", "google", "credentials")
	os.MkdirAll(tokenCacheDir, 0700)

	return filepath.Join(tokenCacheDir, url.QueryEscape(This.Name + ".json")), err
}

func (This *GoogleAnalyticsCollector) tokenFromFile(file string) (*oauth2.Token, error) {
	// tokenFromFile retrieves a Token from a given file path.
	// It returns the retrieved Token and any read error encountered.

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (This *GoogleAnalyticsCollector) saveToken(file string, token *oauth2.Token) error {
	// saveToken uses a file path to create a file and store the
	// token in it.

	logger.I("Saving credential file to: %s", file)

	f, err := os.Create(file)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to cache oauth token: %v", err))
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}