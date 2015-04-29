package pezauth

import (
	"log"

	"github.com/fatih/structs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	EmailFieldName = "email"
)

type (
	//OrgGetHandler - func signature of org get handler
	OrgGetHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//OrgPutHandler - func signature of org put handler
	OrgPutHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
)

type (
	mongoCollection interface {
		Find(query interface{}) *mgo.Query
		Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	}
	persistence interface {
		FindOne(query interface{}, result interface{}) (err error)
		Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	}
	//PivotOrg - struct for pivot org record
	PivotOrg struct {
		Email   string
		OrgName string
	}
	mongoCollectionWrapper struct {
		persistence
		col mongoCollection
	}
	orgController struct {
		Controller
		store persistence
	}
)

//newMongoCollectionWrapper - a wrapper for mongoCollection objects
func newMongoCollectionWrapper(c mongoCollection) persistence {
	return &mongoCollectionWrapper{
		col: c,
	}
}

//FindOne - combining the Find and One calls of a mongo collection object
func (s *mongoCollectionWrapper) FindOne(query interface{}, result interface{}) (err error) {
	err = s.col.Find(query).One(result)
	return
}

//Upsert - allow us to call upsert on mongo collection object
func (s *mongoCollectionWrapper) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return
}

//NewMeController - a controller for me requests
func NewOrgController(c persistence) Controller {
	return &orgController{
		store: c,
	}
}

//Get - get a get handler for org management
func (s *orgController) Get() interface{} {
	var handler OrgGetHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		var (
			result = new(PivotOrg)
		)
		userInfo := GetUserInfo(tokens)
		username := params[UserParam]
		log.Println("getting userInfo: ", userInfo)
		log.Println("result value: ", result)
		err := s.store.FindOne(bson.M{EmailFieldName: username}, result)
		log.Println("response: ", result, err)
		genericResponseFormatter(r, "", structs.Map(result), nil)
	}
	return handler
}

//Put - get a get handler for org management
func (s *orgController) Put() interface{} {
	var handler OrgPutHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		genericResponseFormatter(r, "", nil, nil)
	}
	return handler
}
