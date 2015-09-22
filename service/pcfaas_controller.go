package pezauth

import (
	"log"

	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/pivotal-pez/pezauth/integrations"
	"github.com/go-martini/martini"
)

//NewPcfaasController - a controller for inventory requests
func NewPcfaasController(invClient *integrations.MyInventoryClient) Controller {
	controller := new(pcfaasController)
	controller.inventoryClient = invClient
	return controller
}

// Get - gets a handler for inventory requests
func (s *pcfaasController) Get() interface{} {
	var handler PcfaasGetInventoryHandler = func(log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		statusCode := SuccessStatus
		items, err := s.inventoryClient.GetInventoryItems()
		if err != nil {
			log.Fatalln(err)
			statusCode = ServerErrorStatus
		}
		log.Println("Queried inventory service: ", len(items))
		r.JSON(statusCode, items)
	}
	return handler
}

// Post - Creates a handler for posting to PCFaas to create a new lease.
func (s *pcfaasController) Post() interface{} {
	var handler PcfaasPostInventoryHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		statusCode := SuccessStatus
		userInfo := GetUserInfo(tokens)
		itemID := params[InventoryItemParam]
		emails := userInfo["emails"].([]interface{})
		userEmail := emails[0].(map[string]interface{})["value"]
		lease, err := s.inventoryClient.LeaseInventoryItem(itemID, userEmail.(string), 14)
		if err != nil {
			log.Fatalln(err)
			statusCode = ServerErrorStatus
		}
		log.Println("User info ", userEmail)
		log.Println("Creating a lease for item ", itemID)
		r.JSON(statusCode, lease)
	}
	return handler
}
