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
		Conn redis.Conn
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
)
