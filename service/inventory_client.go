package pezauth

import (
  "net/http"
  "encoding/json"
  "labix.org/v2/mgo/bson"
  "fmt"
  "errors"
  "strings"
  "log"
  "io/ioutil"
)

// NewInventoryServiceClient - creates a new inventory client for communicating
// with the inventory service.
func NewInventoryServiceClient(baseURL string) InventoryServiceClient {
	client := &inventoryServiceClient{serviceBaseURL: baseURL}
  return client
}

type (
  inventoryServiceClient struct {
    serviceBaseURL    string
  }

  //InventoryResponseMessage structures output into a standard format.
  InventoryResponseMessage struct {
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
  // InventoryServiceItem - Represents an inventory item as returned from
  // the remote inventory service.
  InventoryServiceItem struct {
   ID           bson.ObjectId          `bson:"_id,omitempty"`
   SKU          string                 `json:"sku"`
   Tier         int                    `json:"tier"`
   OfferingType string                 `json:"offeringType"`
   Size         string                 `json:"size"`
   Attributes   map[string]interface{} `json:"attributes"`
   Status       string                 `json:"status"`
   LeaseID      string                 `json:"lease_id"`
 }

 // InventoryServiceClient - Represents the set of methods available for the
 // service client.
 InventoryServiceClient interface {
  GetInventoryItems() (result []InventoryItem, err error)
 }
)

func (client *inventoryServiceClient) GetInventoryItems() (result []InventoryItem, err error) {
  effectiveURL := fmt.Sprint(client.serviceBaseURL, "/inventory")
  fmt.Println(effectiveURL)
  r, err := http.Get(effectiveURL)
  if err != nil {
    log.Fatal(err)
    return
  }
  response, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
    return
  }
  r.Body.Close()
  var responseMessage InventoryResponseMessage
  err = json.Unmarshal(response, &responseMessage)
  var inventoryItems []InventoryServiceItem
  b, err := json.Marshal(responseMessage.Data)
  err = json.Unmarshal(b, &inventoryItems)
  result = make([]InventoryItem, len(inventoryItems))
  if strings.ToUpper(responseMessage.Status) != "SUCCESS" {
    err = errors.New(responseMessage.Message)
    return
  }
  for idx, element := range inventoryItems {
    result[idx] = InventoryItem{SKU: element.SKU, Tier: string(element.Tier), OfferingType: element.OfferingType, Status: element.Status, ID: element.ID.Hex()}
  }
  return
}
