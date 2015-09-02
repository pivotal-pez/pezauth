package pezauth

import (
	"log"

	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

//NewMeController - a controller for me requests
func NewPcfaasController() Controller {
	return new(pcfaasController)
}

//Get - get a get handler for authkeyv1
func (s *pcfaasController) Get() interface{} {
	var handler PcfaasGetInventoryHandler = func(log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		userInfo := GetUserInfo(tokens)
		log.Println("getting userInfo: ", userInfo)
    statusCode := SuccessStatus
	  invItem := InventoryItem{SKU:"2C.small",Tier:"2",OfferingType:"C",Size:"small",Status:"available",ID:"abc123guid"}
    invItem2 := InventoryItem{SKU:"2C.small", Tier:"2", OfferingType:"C", Size:"small", Status:"leased", ID:"abc32123guid"}
    items := make([]InventoryItem, 2)
    items[0] = invItem
    items[1] = invItem2
    r.JSON(statusCode, items)
	}
	return handler
}
