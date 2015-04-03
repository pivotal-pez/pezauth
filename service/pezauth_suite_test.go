package pezauth_test

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

func TestPezAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pez Auth Suite")
}

func setVcapApp() {
	os.Setenv("VCAP_APPLICATION", `{  "application_name": "pezauthdev_73b90a93043eb59ee9b3d202dd525f762e865130",  "application_uris": [   "http://localhost:3000"  ],  "application_version": "d744bf29-1465-4634-905d-4fd8a1c19777",  "limits": {   "disk": 1024,   "fds": 16384,   "mem": 1024  },  "name": "pezauthdev_73b90a93043eb59ee9b3d202dd525f762e865130",  "space_id": "49b3e004-702a-4f2c-835c-f25d022882c9",  "space_name": "pez-test",  "uris": [   "http://localhost:3000"  ],  "users": null,  "version": "d744bf29-1465-4634-905d-4fd8a1c19777" }`)
}

func setVcapServ() {
	os.Setenv("VCAP_SERVICES", `{ }`)
}

type mockTokens struct{}

func (s *mockTokens) Access() (r string) {
	return
}

func (s *mockTokens) Refresh() (r string) {
	return
}
func (s *mockTokens) Expired() (r bool) {
	return
}
func (s *mockTokens) ExpiryTime() (r time.Time) {
	return
}

type mockResponseWriter struct {
	StatusCode int
}

func (s *mockResponseWriter) WriteHeader(i int) {
	s.StatusCode = i
}

func (s *mockResponseWriter) Header() (r http.Header) {
	return
}

func (s *mockResponseWriter) Write(x []byte) (a int, b error) {
	return
}
