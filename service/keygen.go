package pezauth

import "fmt"

//KeyGenerator - interface to work with apikeys
type KeyGenerator interface {
	Get(user string) (string, error)
	Create(user string) error
	Delete(user string) error
}

//Doer - interface to make a call to persistence store
type Doer interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

//NewKeyGen - create a new implementation of a KeyGenerator interface
func NewKeyGen(doer Doer, guid GUIDMaker) KeyGenerator {
	return &KeyGen{
		store:     doer,
		guidMaker: guid,
	}
}

//KeyGen - and implementation of the KeyGenerator interface
type KeyGen struct {
	store     Doer
	guidMaker GUIDMaker
}

func parseScanResponse(r interface{}) (res string) {
	resArr := r.([]interface{})
	responseArrayIdx := 1
	first := 0

	if len(resArr) > responseArrayIdx {

		if arr := resArr[responseArrayIdx]; len(arr.([]interface{})) > first {
			res = arr.([]interface{})[first].(string)
		}
	}
	return
}

//Get - gets a key for a user
func (s *KeyGen) Get(user string) (res string, err error) {
	var r interface{}
	search := fmt.Sprintf("%s:*", user)

	if r, err = s.store.Do("SCAN", 0, "MATCH", search); r != nil {
		res = parseScanResponse(r)
	}
	return
}

func createHash(user, guid string) (hash string) {
	hash = fmt.Sprintf("%s:%s", user, guid)
	return
}

//Create - creates a new key for a user
func (s *KeyGen) Create(user string) (err error) {
	guid := s.guidMaker.Create()
	hash := createHash(user, guid)
	_, err = s.store.Do("HMSET", hash)
	return
}

//Delete - deletes a key for a user
func (s *KeyGen) Delete(user string) (err error) {
	var hash string

	if hash, err = s.Get(user); err == nil {
		_, err = s.store.Do("DEL", hash)
	}
	return
}
