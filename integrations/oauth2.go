package integrations

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

//New - create a new oauth2 integration wrapper
func (s *MyOAuth2) New(appEnv *cfenv.App) *MyOAuth2 {
	oauth2ServiceName := os.Getenv("OAUTH2_SERVICE_NAME")
	clientIdName := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecretName := os.Getenv("OAUTH2_CLIENT_SECRET")
	oauth2Service, err := appEnv.Services.WithName(oauth2ServiceName)

	if err != nil {
		panic(fmt.Sprintf("oauth2 client service name error: %s", err.Error()))
	}
	s.ID = oauth2Service.Credentials[clientIdName]
	s.Secret = oauth2Service.Credentials[clientSecretName]
	return s
}
