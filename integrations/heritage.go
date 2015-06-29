package integrations

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

//New - create a new heritage foundation integration wrapper
func (s *MyHeritage) New(appEnv *cfenv.App) *MyHeritage {
	heritageAdminServiceName := os.Getenv("UPS_PEZ_HERITAGE_ADMIN_NAME")
	heritageLoginTargetName := os.Getenv("HERITAGE_LOGIN_TARGET_NAME")
	heritageLoginUserName := os.Getenv("HERITAGE_LOGIN_USER_NAME")
	heritageLoginPassName := os.Getenv("HERITAGE_LOGIN_PASS_NAME")
	heritageCCTargetName := os.Getenv("HERITAGE_CC_TARGET_NAME")
	heritageAdminService, err := appEnv.Services.WithName(heritageAdminServiceName)

	if err != nil {
		panic(fmt.Sprintf("heritage service name error: %s", err.Error()))
	}
	s.LoginTarget = heritageAdminService.Credentials[heritageLoginTargetName].(string)
	s.LoginUser = heritageAdminService.Credentials[heritageLoginUserName].(string)
	s.LoginPass = heritageAdminService.Credentials[heritageLoginPassName].(string)
	s.CCTarget = heritageAdminService.Credentials[heritageCCTargetName].(string)
	return s
}
