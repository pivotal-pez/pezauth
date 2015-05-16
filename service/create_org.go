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
	ListUsersEndpoint = "/Users"
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
	//UserAPIResponse - the user api response object
	UserAPIResponse struct {
		Resources []map[string]interface{}
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
		s.authRequestor(s.authClient.CCTarget(), "GET", "", GetAPIInfo, SuccessStatus, func(res *http.Response) {
			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(b, &s.apiInfo)
			s.log.Println(s.apiInfo)
			s.log.Println(string(b[:]))

		}, func(res *http.Response, e error) {
			b, _ := ioutil.ReadAll(res.Body)
			s.log.Println("error: ", e, string(b[:]))
		})
	} else {
		s.log.Println("wtf is going on here")
	}
}

func (s *orgManager) getAuthEnpoint() (endpoint string) {
	s.setApiInfo()

	switch authEndpoint := s.apiInfo["authorization_endpoint"].(type) {
	case string:
		endpoint = authEndpoint
	}
	return
}

func (s *orgManager) getUAAEnpoint() (endpoint string) {
	s.setApiInfo()

	switch tokenEndpoint := s.apiInfo["token_endpoint"].(type) {
	case string:
		endpoint = tokenEndpoint
	}
	return
}

func (s *orgManager) getUserGUID() (guid string, err error) {
	var (
		userResponse UserAPIResponse
		data         = map[string]string{
			"attributes": "id,userName",
			//"filter":     fmt.Sprintf("userName Eq %s", s.username),
		}
	)

	s.authRequestor(s.getUAAEnpoint(), "GET", data, ListUsersEndpoint, SuccessStatus, func(res *http.Response) {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(b, &userResponse)

		for _, resource := range userResponse.Resources {

			switch id := resource["id"].(type) {
			case string:

				if resource["userName"] == s.username {
					guid = id
					s.log.Println("we have a user guid!", guid)
				}
			}
		}

	}, func(res *http.Response, e error) {
		b, _ := ioutil.ReadAll(res.Body)
		s.log.Println("call for user guid failed :(", e, string(b[:]))

	})
	return
}

func (s *orgManager) addRoles(orgGUID string) (err error) {
	var (
		userGUID string
	)
	s.log.Println("creating a role for orgguid: ", orgGUID)

	if userGUID, err = s.getUserGUID(); err == nil {
		managerPath := fmt.Sprintf("/v2/organizations/%s/managers/%s", orgGUID, userGUID)
		usersPath := fmt.Sprintf("/v2/organizations/%s/users/%s", orgGUID, userGUID)
		s.addRoleFromPath(managerPath)
		s.addRoleFromPath(usersPath)
	}
	return
}

func (s *orgManager) addRoleFromPath(rolePath string) {
	s.authRequestor(s.authClient.CCTarget(), "PUT", "", rolePath, OrgCreateSuccessStatusCode, func(res *http.Response) {
		s.log.Println("we have a role!", rolePath)

	}, func(res *http.Response, e error) {
		b, _ := ioutil.ReadAll(res.Body)
		s.log.Println("call for role failed :(", rolePath, e, res.StatusCode, string(b[:]))

	})
}

func (s *orgManager) Create() (record *PivotOrg, err error) {
	var (
		data = fmt.Sprintf(`{"name":"%s"}`, getOrgNameFromEmail(s.username))
	)

	s.authRequestor(s.authClient.CCTarget(), "POST", data, OrgCreateEndpoint, OrgCreateSuccessStatusCode, func(res *http.Response) {
		s.log.Println("we created the org successfully")
		guid := s.parseOrgGUID(res)

		if err = s.addRoles(guid); err == nil {
			record, err = s.upsert(guid)
		}
	}, func(res *http.Response, e error) {
		s.log.Println("call to create org api failed")
		record = new(PivotOrg)
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

func (s *orgManager) upsert(orgGUID string) (record *PivotOrg, err error) {
	orgname := getOrgNameFromEmail(s.username)
	record = &PivotOrg{
		Email:   s.username,
		OrgName: orgname,
		OrgGUID: orgGUID,
	}
	err = s.store.Upsert(bson.M{EmailFieldName: s.username}, record)
	return
}

func (s *orgManager) authRequestor(url string, verb string, data interface{}, path string, successCode int, callback func(*http.Response), callbackFail func(*http.Response, error)) {
	var (
		req *http.Request
		res *http.Response
		err error
	)
	s.log.Println("making rest call to: ", url, "-", verb, "-", data, "-", path)

	if req, err = s.authClient.CreateAuthRequest(verb, url, path, data); err == nil {
		s.log.Println("we created the decorated request")

		if res, err = s.authClient.HttpClient().Do(req); res.StatusCode == successCode && err == nil {
			s.log.Println("we are now going to execute the success callback")
			defer res.Body.Close()
			callback(res)

		} else {
			s.log.Println("we are now going to execute the failure callback")
			callbackFail(res, err)
		}
	}
}

func getOrgNameFromEmail(email string) (orgName string) {
	username := strings.Split(email, "@")[0]
	orgName = fmt.Sprintf("pivot-%s", username)
	return
}
