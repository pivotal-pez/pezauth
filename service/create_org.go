package pezauth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/martini-contrib/oauth2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrOrgCreateAPICallFailure = errors.New("failed to create org on api call")
)

type (
	orgManager struct {
		username   string
		log        *log.Logger
		tokens     oauth2.Tokens
		store      persistence
		authClient authRequestCreator
	}
)

func newOrg(username string, log *log.Logger, tokens oauth2.Tokens, store persistence, authClient authRequestCreator) *orgManager {
	return &orgManager{
		username:   username,
		log:        log,
		tokens:     tokens,
		store:      store,
		authClient: authClient,
	}
}

func (s *orgManager) Show() (result *PivotOrg, err error) {
	result = new(PivotOrg)
	userInfo := GetUserInfo(s.tokens)
	NewUserMatch().UserInfo(userInfo).UserName(s.username).OnSuccess(func() {
		log.Println("getting userInfo: ", userInfo)
		log.Println("result value: ", result)
		err = s.store.FindOne(bson.M{EmailFieldName: s.username}, result)
		log.Println("response: ", result, err)
	}).OnFailure(func() {
		log.Println(ErrCantCallAcrossUsers.Error())
		err = ErrCantCallAcrossUsers
	}).Run()
	return
}

func (s *orgManager) Create() (record *PivotOrg, err error) {
	var (
		res  *http.Response
		req  *http.Request
		data = map[string]string{
			"name": getOrgNameFromEmail(s.username),
		}
	)

	if req, err = s.authClient.CreateAuthRequest("POST", s.authClient.CCTarget(), "/v2/organizations", data); err == nil {

		if res, err = s.authClient.HttpClient().Do(req); res.StatusCode == 201 && err == nil {
			record, err = s.upsert()

		} else {
			record = new(PivotOrg)
			err = ErrOrgCreateAPICallFailure
		}
	}
	return
}

func (s *orgManager) upsert() (record *PivotOrg, err error) {
	orgname := getOrgNameFromEmail(s.username)
	record = &PivotOrg{
		Email:   s.username,
		OrgName: orgname,
	}
	s.store.Upsert(bson.M{EmailFieldName: s.username}, record)
	return
}

func getOrgNameFromEmail(email string) (orgName string) {
	username := strings.Split(email, "@")[0]
	orgName = fmt.Sprintf("pivot-%s", username)
	return
}
