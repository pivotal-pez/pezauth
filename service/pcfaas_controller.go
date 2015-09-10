package pezauth

import (
	"log"

	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/pivotal-pez/pezauth/integrations"
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
