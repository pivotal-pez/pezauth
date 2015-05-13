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
	//ListUsersEndpoint - get a list of all users in paas
	ListUsersEndpoint = "/v2/users"
	//GetApiInfo - the endpoint to grab api info data
	GetAPIInfo = "/v2/info"
)

var (
	//ErrOrgCreateAPICallFailure - error for failed call to create org endpoint
	ErrOrgCreateAPICallFailure = errors.New("failed to create org on api call")
)

type (
	orgManager struct {
		username   string
		userGUID   string
		log        *log.Logger
		tokens     oauth2.Tokens
		store      persistence
		authClient authRequestCreator
		apiInfo    map[string]interface{}
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
	//APIResponseList - a list of resources or apiresponse objects
	APIResponseList struct {
		Resources []APIResponse `json:"resources"`
	}
)

func newOrg(username string, log *log.Logger, tokens oauth2.Tokens, store persistence, authClient authRequestCreator) *orgManager {
	s := &orgManager{
		username:   username,
		log:        log,
		tokens:     tokens,
		store:      store,
		authClient: authClient,
	}
	return s
}

func (s *orgManager) Show() (result *PivotOrg, err error) {
	result = new(PivotOrg)
	userInfo := GetUserInfo(s.tokens)
	NewUserMatch().UserInfo(userInfo).UserName(s.username).OnSuccess(func() {
		s.log.Println("getting userInfo: ", userInfo)
		s.log.Println("result value: ", result)
		err = s.store.FindOne(bson.M{EmailFieldName: s.username}, result)
		s.log.Println("response: ", result, err)
	}).OnFailure(func() {
		s.log.Println(ErrCantCallAcrossUsers.Error())
		err = ErrCantCallAcrossUsers
	}).Run()
	return
}

func (s *orgManager) setApiInfo() {
	if s.apiInfo == nil {
		s.authRequestor("GET", nil, GetAPIInfo, SuccessStatus, func(res *http.Response) {
			b, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(b, s.apiInfo)
		}, func(res *http.Response, e error) {
			s.log.Println("error: ", e)
		})
	}
}

func (s *orgManager) getAuthEnpoint() string {
	s.setApiInfo()
	return s.apiInfo["authorization_endpoint"].(string)
}

func (s *orgManager) addRoles() (err error) {
	s.log.Println("we still need to implement role creation")
	return
}

func (s *orgManager) Create() (record *PivotOrg, err error) {
	var (
		data = map[string]string{"abc": fmt.Sprintf(`{"name":"%s"}`, getOrgNameFromEmail(s.username))}
	)

	s.authRequestor("POST", data, OrgCreateEndpoint, OrgCreateSuccessStatusCode, func(res *http.Response) {
		s.log.Println("we created the org successfully")

		if err = s.addRoles(); err == nil {
			record, err = s.upsert(res)
		}
	}, func(res *http.Response, e error) {
		s.log.Println("call to create org api failed")
		record = new(PivotOrg)
		s.log.Println("we are seeing this type of error: ", e)
		b, _ := ioutil.ReadAll(res.Body)
		s.log.Println("we are seeing this type of response: ", string(b[:]))
		err = ErrOrgCreateAPICallFailure
	})
	return
}

func (s *orgManager) parseOrgGUID(res *http.Response) (guid string) {
	apiResponse := new(APIResponse)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, apiResponse)
	guid = apiResponse.Metadata.GUID
	return
}

func (s *orgManager) upsert(res *http.Response) (record *PivotOrg, err error) {
	orgname := getOrgNameFromEmail(s.username)
	guid := s.parseOrgGUID(res)
	record = &PivotOrg{
		Email:   s.username,
		OrgName: orgname,
		OrgGUID: guid,
	}
	err = s.store.Upsert(bson.M{EmailFieldName: s.username}, record)
	return
}

func (s *orgManager) authRequestor(verb string, data map[string]string, path string, successCode int, callback func(*http.Response), callbackFail func(*http.Response, error)) {
	var (
		req *http.Request
		res *http.Response
		err error
	)
	s.authClient.ParseDataAsString(true)

	if req, err = s.authClient.CreateAuthRequest(verb, s.authClient.CCTarget(), path, data); err == nil {

		if res, err = s.authClient.HttpClient().Do(req); res.StatusCode == successCode && err == nil {
			defer res.Body.Close()
			callback(res)

		} else {
			callbackFail(res, err)
		}
	}
}

func getOrgNameFromEmail(email string) (orgName string) {
	username := strings.Split(email, "@")[0]
	orgName = fmt.Sprintf("pivot-%s", username)
	return
}
