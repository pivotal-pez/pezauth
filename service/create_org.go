package pezauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/martini-contrib/oauth2"
	"gopkg.in/mgo.v2/bson"
)

const (
	//OrgCreateSuccessStatusCode - success status code from a call to the org create cc endpoint
	OrgCreateSuccessStatusCode = 201
	//OrgCreateEndpoint - the endpoint to hit for org creation in the cc api
	OrgCreateEndpoint = "/v2/organizations"
)

var (
	//ErrOrgCreateAPICallFailure - error for failed call to create org endpoint
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
	//APIResponse - cc http response object
	APIResponse struct {
		Metadata APIMetadata            `json:"metadata"`
		Entity   map[string]interface{} `json:"entity"`
	}
	//APIMetadata = cc http response metadata
	APIMetadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
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

	if req, err = s.authClient.CreateAuthRequest("POST", s.authClient.CCTarget(), OrgCreateEndpoint, data); err == nil {

		if res, err = s.authClient.HttpClient().Do(req); res.StatusCode == OrgCreateSuccessStatusCode && err == nil {
			defer res.Body.Close()
			record, err = s.upsert(res)

		} else {
			record = new(PivotOrg)
			err = ErrOrgCreateAPICallFailure
		}
	}
	return
}

func (s *orgManager) parseGUID(res *http.Response) (guid string) {
	apiResponse := new(APIResponse)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, apiResponse)
	guid = apiResponse.Metadata.GUID
	return
}

func (s *orgManager) upsert(res *http.Response) (record *PivotOrg, err error) {
	orgname := getOrgNameFromEmail(s.username)
	guid := s.parseGUID(res)
	record = &PivotOrg{
		Email:   s.username,
		OrgName: orgname,
		OrgGUID: guid,
	}
	s.store.Upsert(bson.M{EmailFieldName: s.username}, record)
	return
}

func getOrgNameFromEmail(email string) (orgName string) {
	username := strings.Split(email, "@")[0]
	orgName = fmt.Sprintf("pivot-%s", username)
	return
}
