package pezauth_test

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/xchapter7x/cloudcontroller-client"
)

const (
	QueryUserUAAResponseBody = `{
  "resources": [
    {
      "id": "123456"
    }
  ],
  "startIndex": 1,
  "itemsPerPage": 100,
  "totalResults": 1,
  "schemas":["urn:scim:schemas:core:1.0"]
}`
	CreateOrgResponseBody = `{
  "metadata": {
    "guid": "8b0939ca-3a69-40c4-aa12-f771f3e1cf3e",
    "url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e",
    "created_at": "2015-06-16T23:40:50Z",
    "updated_at": null
  },
  "entity": {
    "name": "my-org-name",
    "billing_enabled": false,
    "quota_definition_guid": "5504af1a-ff01-44d1-9bd2-d2d799bdd58a",
    "status": "active",
    "quota_definition_url": "/v2/quota_definitions/5504af1a-ff01-44d1-9bd2-d2d799bdd58a",
    "spaces_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/spaces",
    "domains_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/domains",
    "private_domains_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/private_domains",
    "users_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/users",
    "managers_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/managers",
    "billing_managers_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/billing_managers",
    "auditors_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/auditors",
    "app_events_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/app_events",
    "space_quota_definitions_url": "/v2/organizations/8b0939ca-3a69-40c4-aa12-f771f3e1cf3e/space_quota_definitions"
  }
}`
	CreateSpaceResponseBody = `{
  "metadata": {
    "guid": "6c1fe603-71f1-4f37-93fb-04b674d9bebc",
    "url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc",
    "created_at": "2015-06-16T23:41:05Z",
    "updated_at": null
  },
  "entity": {
    "name": "development",
    "organization_guid": "6a6e43ef-7337-47f7-b3ad-6a6cee3b5bce",
    "space_quota_definition_guid": null,
    "allow_ssh": true,
    "organization_url": "/v2/organizations/6a6e43ef-7337-47f7-b3ad-6a6cee3b5bce",
    "developers_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/developers",
    "managers_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/managers",
    "auditors_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/auditors",
    "apps_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/apps",
    "routes_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/routes",
    "domains_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/domains",
    "service_instances_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/service_instances",
    "app_events_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/app_events",
    "events_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/events",
    "security_groups_url": "/v2/spaces/6c1fe603-71f1-4f37-93fb-04b674d9bebc/security_groups"
  }
}`
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

func NewMockDoer() ccclient.ClientDoer {
	return &mockOrgChainDoer{
		CallChainCounter: 0,
		CallChainResponses: []callChainResponse{
			ccCallChain(200, "Login"),
			ccCallChain(200, QueryUserUAAResponseBody),
			ccCallChain(201, CreateOrgResponseBody),
			ccCallChain(201, "added org manager role!"),
			ccCallChain(201, "added org user role!"),
			ccCallChain(201, CreateSpaceResponseBody),
			ccCallChain(201, "added space manager role!"),
			ccCallChain(201, "added space user role!"),
		},
	}
}

func ccCallChain(statusCode int, body string) callChainResponse {
	return callChainResponse{
		Res: &http.Response{
			StatusCode: statusCode,
			Body:       &ClosingBuffer{bytes.NewBufferString(body)},
		},
		Err: nil,
	}
}

type mockOrgChainDoer struct {
	CallChainCounter   int
	CallChainResponses []callChainResponse
}

func (s *mockOrgChainDoer) Do(*http.Request) (*http.Response, error) {
	res := s.CallChainResponses[s.CallChainCounter]
	s.CallChainCounter++
	return res.Res, res.Err
}

type callChainResponse struct {
	Res *http.Response
	Err error
}

type mockAuthRequestCreator struct {
	MyMockDoer ccclient.ClientDoer
}

func (s *mockAuthRequestCreator) CreateAuthRequest(verb, requestURL, path string, args interface{}) (req *http.Request, err error) {
	req = new(http.Request)
	req.URL = new(url.URL)
	return
}
func (s *mockAuthRequestCreator) CCTarget() (r string) {
	return
}
func (s *mockAuthRequestCreator) HttpClient() (do ccclient.ClientDoer) {
	do = s.MyMockDoer
	return
}
func (s *mockAuthRequestCreator) Login() (clnt *ccclient.Client, err error) {
	return
}

type mockOrgPersistence struct {
	CalledUpsert int
}

func (s *mockOrgPersistence) Remove(selector interface{}) (err error) {
	return
}
func (s *mockOrgPersistence) FindOne(query interface{}, result interface{}) (err error) {
	return
}
func (s *mockOrgPersistence) Upsert(selector interface{}, update interface{}) (err error) {
	s.CalledUpsert++
	return
}
