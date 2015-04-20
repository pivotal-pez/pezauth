package pezauth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
)

var (
	ErrUnparsableHash   = errors.New("Could not parse the hash or hash was nil")
	ErrEmptyKeyResponse = errors.New("The key could not be found or was not valid")
)

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

func parseKeysResponse(r interface{}) (key, username, hash string, err error) {

	if resArr := r.([]interface{}); len(resArr) > 0 {
		ba := resArr[0].([]byte)
		hash = string(ba[:])
		key, username, err = hashSplit(hash)

	} else {
		err = ErrEmptyKeyResponse
	}
	return
}

//Get - gets a key for a user
func (s *KeyGen) Get(user string) (res string, err error) {
	var r interface{}
	search := fmt.Sprintf("%s:*", user)

	if r, err = s.store.Do("KEYS", search); r != nil && err == nil {
		res, _, _, err = parseKeysResponse(r)
	}
	return
}

func (s *KeyGen) getHash(user string) (hash string, err error) {
	var r interface{}
	search := fmt.Sprintf("%s:*", user)

	if r, err = s.store.Do("KEYS", search); r != nil && err == nil {
		_, _, hash, err = parseKeysResponse(r)
	}
	return
}

func hashSplit(hash string) (key, username string, err error) {
	usernameIndex := 0
	keyIndex := 1
	hashSplitArrayLen := 2

	if splitHash := strings.Split(hash, ":"); len(splitHash) == hashSplitArrayLen {
		key = splitHash[keyIndex]
		username = splitHash[usernameIndex]

	} else {
		err = ErrUnparsableHash
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
	row := map[string]string{"active": "true"}
	_, err = s.store.Do("HMSET", redis.Args{hash}.AddFlat(row)...)
	return
}

//Delete - deletes a key for a user
func (s *KeyGen) Delete(user string) (err error) {
	var apikey string

	if apikey, err = s.Get(user); err == nil {
		fmt.Println("we should now be deleting:", apikey)
		_, err = s.store.Do("DEL", createHash(user, apikey))
		fmt.Println("the error?:", err)
	}
	return
}
