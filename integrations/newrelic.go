package integrations

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

//New - create a new newrelic integration wrapper
func (s *MyNewRelic) New(appEnv *cfenv.App) *MyNewRelic {
	serviceName := os.Getenv("NEWRELIC_SERVICE_NAME")
	keyName := os.Getenv("NEWRELIC_KEY_NAME")
	appName := os.Getenv("NEWRELIC_APP_NAME")
	service, err := appEnv.Services.WithName(serviceName)

	if err != nil {
		panic(fmt.Sprintf("new relic service name error: %s", err.Error()))
	}
	s.Key = service.Credentials[keyName].(string)
	s.App = service.Credentials[appName].(string)
	return s
}
