package pezauth

type KeyGenerator interface {
	Get(user string) (string, error)
	Create(user string) error
	Delete(user string) error
}

type Doer interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

func NewKeyGen(doer Doer, guid GUIDMaker) KeyGenerator {
	return &KeyGen{
		store:     doer,
		guidMaker: guid,
	}
}

type KeyGen struct {
	store     Doer
	guidMaker GUIDMaker
}

func (s *KeyGen) Get(user string) (string, error) {
	r, err := s.store.Do("GET", user)
	return r.(string), err
}

func (s *KeyGen) Create(user string) (err error) {
	guid := s.guidMaker.Create()
	_, err = s.store.Do("SET", user, guid)
	return
}

func (s *KeyGen) Delete(user string) (err error) {
	_, err = s.store.Do("DEL", user)
	return
}
