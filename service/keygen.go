package pezauth

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

//Get - gets a key for a user
func (s *KeyGen) Get(user string) (string, error) {
	r, err := s.store.Do("GET", user)
	return r.(string), err
}

//Create - creates a new key for a user
func (s *KeyGen) Create(user string) (err error) {
	guid := s.guidMaker.Create()
	_, err = s.store.Do("SET", user, guid)
	return
}

//Delete - deletes a key for a user
func (s *KeyGen) Delete(user string) (err error) {
	_, err = s.store.Do("DEL", user)
	return
}
