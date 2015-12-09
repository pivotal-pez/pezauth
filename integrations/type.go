package integrations

import (
	"github.com/garyburd/redigo/redis"
	"labix.org/v2/mgo"
)

type (
	//MyOAuth2 - integration wrapper for oauth2
	MyOAuth2 struct {
		ID     string
		Secret string
	}
	//MyRedis - integration wrapper for redis
	MyRedis struct {
		URI  string
		Pass string
		Pool *redis.Pool
	}
	//MyMongo - integration wrapper for mongodb
	MyMongo struct {
		Col                *mgo.Collection
		Session            *mgo.Session
		mongoConnectionURI string
		mongoDBName        string
		mongoCollName      string
	}
	//MyHeritage - integration wrapper for connections to heritage
	MyHeritage struct {
		LoginTarget string
		LoginUser   string
		LoginPass   string
		CCTarget    string
	}
	//MyNewRelic - integration wrapper for connections to newrelic
	MyNewRelic struct {
		Key string
		App string
	}
	//MyInventoryClient - integration wrapper for interacting with inventory svc
	MyInventoryClient struct {
		ServiceBaseURL string
		Enabled        bool
	}

	// InventoryItem - entity from inventory query, includes lease status
	InventoryItem struct {
		SKU          string         `json:"sku"`
		Tier         int            `json:"tier"`
		OfferingType string         `json:"offeringType"`
		Size         string         `json:"size"`
		Status       string         `json:"status"`
		ID           string         `json:"id"`
		CurrentLease InventoryLease `json:"currentLease"`
	}

	// InventoryLease - represents information about an active lease of an inventory item.
	InventoryLease struct {
		DaysUntilExpires int    `json:"daysUntilExpires"`
		Username         string `json:"userName"`
	}
)
