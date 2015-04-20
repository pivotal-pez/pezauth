package pezauth

import "errors"

var (
	ErrNotValidActionForUser = errors.New("not a valid user to perform this action")
)

func NewUserMatch() *UserMatch {
	return new(UserMatch)
}

type UserMatch struct {
	userInfo    map[string]interface{}
	username    string
	successFunc func()
	failFunc    func()
}

func (s *UserMatch) UserInfo(userInfo map[string]interface{}) *UserMatch {
	s.userInfo = userInfo
	return s
}

func (s *UserMatch) UserName(username string) *UserMatch {
	s.username = username
	return s
}

func (s *UserMatch) OnSuccess(successFunc func()) *UserMatch {
	s.successFunc = successFunc
	return s
}

func (s *UserMatch) OnFailure(failFunc func()) *UserMatch {
	s.failFunc = failFunc
	return s
}

func (s *UserMatch) Run() (err error) {
	var hasValidEmail = false

	for _, email := range s.userInfo["emails"].([]interface{}) {

		if email.(map[string]interface{})["value"].(string) == s.username {
			hasValidEmail = true
			s.successFunc()
		}
	}

	if !hasValidEmail {
		s.failFunc()
		err = ErrNotValidActionForUser
	}
	return
}
