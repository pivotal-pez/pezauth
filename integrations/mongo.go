package integrations

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotal-pez/pezdispenser/service"
	"labix.org/v2/mgo"
)

//New - create a new mongo integration wrapper
func (s *MyMongo) New(appEnv *cfenv.App) *MyMongo {
	mongoServiceName := os.Getenv("MONGO_SERVICE_NAME")
	mongoURIName := os.Getenv("MONGO_URI_NAME")
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	s.mongoCollName = os.Getenv("MONGO_COLLECTION_NAME")
	mongoService, err := appEnv.Services.WithName(mongoServiceName)

	if err != nil {
		panic(fmt.Sprintf("mongodb service name error: %s", err.Error()))
	}
	s.mongoConnectionURI = mongoService.Credentials[mongoURIName].(string)
	s.mongoDBName = mongoService.Credentials[mongoDBName].(string)
	s.connect()
	return s
}

func (s *MyMongo) connect() {
	var err error

	if s.Session, err = mgo.Dial(s.mongoConnectionURI); err != nil {
		panic(fmt.Sprintf("mongodb dial error: %s", err.Error()))
	}
	s.Session.SetMode(mgo.Monotonic, true)
	s.Col = s.Session.DB(s.mongoDBName).C(s.mongoCollName)
}

//Collection - this allows us to get a mongo collection with a new session wrapped as a pezdispenser.Persistence interface implementation
func (s *MyMongo) Collection() pezdispenser.Persistence {
	sess := s.Session.Copy()
	return pezdispenser.NewMongoCollectionWrapper(sess.DB(s.mongoDBName).C(s.mongoCollName))
}
