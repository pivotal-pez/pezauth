package pezauth_test

import (
	"log"
//	"os"
	"fmt"
   "net/http"
	 "net/http/httptest"
	 "strings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezauth/integrations"
)

func makeServer(payload string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")

			fmt.Fprintln(w, payload)
			}))
}

func makeLeaseAvailableServer(payloadFirst string, payloadSecond string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")

			if strings.Contains(r.RequestURI, "lease") {
				fmt.Fprintln(w, payloadSecond)
			} else {
				fmt.Fprintln(w, payloadFirst)
			}
		}))
}

var _ = Describe("Inventory Service Client", func() {
	Context("calling .GetInventoryItems()", func() {

		It("should obtain and properly unmarshal results with data", func() {
			server := makeServer(TwoItemsSample)
			defer server.Close()
			rootURL := server.URL

			invClient := new(MyInventoryClient).NewWithURL(rootURL)
      results, err := invClient.GetInventoryItems()
			Ω(err).Should(BeNil())
      Ω(len(results)).Should(Equal(2))
			Ω(results[0].ID).Should(Equal("55f077966696f40147000001"))
		})

		It("should properly handle a failure in server response message", func() {
			server := makeServer(NoInventoryDataSample)
			defer server.Close()
			rootURL := server.URL
			invClient := new(MyInventoryClient).NewWithURL(rootURL)
	    results, err := invClient.GetInventoryItems()
			Ω(len(results)).Should(Equal(0))
			Ω(err.Error()).Should(Equal("Things went horribly wrong"))
		})

		It("should query for a lease when appropriate", func() {
				server := makeLeaseAvailableServer(LeasedItemsSample, IndividualLease)
				defer server.Close()
				rootURL := server.URL
				invClient := new(MyInventoryClient).NewWithURL(rootURL)
				results, err := invClient.GetInventoryItems()
				Ω(len(results)).Should(Equal(2))
				Ω(err).Should(BeNil())
				log.Printf("Detected user owning lease %s", results[1].CurrentLease.Username)
				Ω(results[1].CurrentLease.Username).Should(Equal("dnem"))
		})
	})
})
