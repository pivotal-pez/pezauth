package integrations

import (
  "net/http"
  "encoding/json"
  "fmt"
  "errors"
  "log"
  "io/ioutil"
  "labix.org/v2/mgo/bson"
  "os"
  "strings"
  "crypto/tls"
  cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// New - create an inventory service wrapper.
func (client *MyInventoryClient) New(appEnv *cfenv.App) *MyInventoryClient {

  if strings.ToUpper(os.Getenv("DISPLAY_NEW_SERVICES")) == "YES" {
    client.Enabled = true
    client.ServiceBaseURL = getServiceBinding("inventory-service", "target-url", appEnv)
  } else {
    client.Enabled = false
  }
  return client
}

// NewWithURL - create client directly with URL. Used for testing.
func (client *MyInventoryClient) NewWithURL(url string) *MyInventoryClient {
  client.Enabled = true
  client.ServiceBaseURL = url
  return client
}

// GetInventoryItems - query the inventory items from the inventory service
func (client *MyInventoryClient) GetInventoryItems() (result []InventoryItem, err error) {
  if !client.Enabled {
    return
  }
  effectiveURL := fmt.Sprint(client.ServiceBaseURL, "/inventory")
  // Currently allowing invalid certs to talk to inventory service.
  // required to talk to inventory service locally or in dev.
  // TODO make this a configurable toggle so we can enforce it in prod.
  tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  httpclient := &http.Client{Transport: tr}
  r, err := httpclient.Get(effectiveURL)
  if err != nil {
    log.Println(err)
    return
  }
  response, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    return
  }
  r.Body.Close()
  var responseMessage inventoryResponseMessage
  err = json.Unmarshal(response, &responseMessage)
  var inventoryItems []inventoryServiceItem
  b, err := json.Marshal(responseMessage.Data)
  err = json.Unmarshal(b, &inventoryItems)
  result = make([]InventoryItem, len(inventoryItems))
  if strings.ToUpper(responseMessage.Status) != "SUCCESS" {
    err = errors.New(responseMessage.Message)
    return
  }
  for idx, element := range inventoryItems {
    result[idx] = InventoryItem{SKU: element.SKU, Tier: string(element.Tier), OfferingType: element.OfferingType, Status: element.Status, ID: element.ID.Hex()}
    // TODO - fetch the lease information for this inventory item if there is a lease...
  }
  return
}



func (client *MyInventoryClient) String() string {
  if client.Enabled {
    return fmt.Sprintf("{InventoryClient pointing at %s }", client.ServiceBaseURL)
  }
  return fmt.Sprintf("{InventoryClient DISABLED}")
}

func getServiceBinding(serviceName string, serviceURIName string, appEnv *cfenv.App) (serviceURI string) {

	if service, err := appEnv.Services.WithName(serviceName); err == nil {
		if serviceURI = service.Credentials[serviceURIName].(string); serviceURI == "" {
			panic(fmt.Sprintf("we pulled an empty connection string %s from %v - %v", serviceURI, service, service.Credentials))
		}

	} else {
		panic(fmt.Sprint("Experienced an error trying to grab service binding information:", err.Error()))
	}
	return
}

// Private types
// ------------------------------------------------

type inventoryResponseMessage struct {
    //Status returns a string indicating [success|error|fail]
    Status string `json:"status"`
    //Data holds the payload of the response
    Data interface{} `json:"data,omitempty"`
    //Message contains the nature of an error
    Message string `json:"message,omitempty"`
    //Meta contains information about the data and the current request
    Meta map[string]interface{} `json:"_metaData,omitempty"`
    //Links contains [prev|next] links for paginated responses
    Links map[string]interface{} `json:"_links,omitempty"`
}

type inventoryServiceItem struct {
 ID           bson.ObjectId          `bson:"_id,omitempty"`
 SKU          string                 `json:"sku"`
 Tier         int                    `json:"tier"`
 OfferingType string                 `json:"offering_type"`
 Size         string                 `json:"size"`
 Attributes   map[string]interface{} `json:"attributes"`
 Status       string                 `json:"status"`
 LeaseID      string                 `json:"lease_id"`
}
