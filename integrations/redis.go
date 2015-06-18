package integrations

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
)

//New - create a new redis integration wrapper
func (s *MyRedis) New(appEnv *cfenv.App) *MyRedis {
	redisName := os.Getenv("REDIS_SERVICE_NAME")
	redisHost := os.Getenv("REDIS_HOSTNAME_NAME")
	redisPass := os.Getenv("REDIS_PASSWORD_NAME")
	redisPort := os.Getenv("REDIS_PORT_NAME")
	redisService, err := appEnv.Services.WithName(redisName)
	s.Pass = redisService.Credentials[redisPass]
	s.URI = fmt.Sprintf("%s:%s", redisService.Credentials[redisHost], redisService.Credentials[redisPort])

	if err != nil {
		panic(fmt.Sprintf("redis service name error: %s", err.Error()))
	}
	s.connect()
	defer func() { go s.autoReconnect() }()
	return s
}

func (s *MyRedis) connect() {
	var err error

	if s.Conn, err = redis.Dial("tcp", s.URI); err != nil {
		panic(fmt.Sprintf("redis dial error: %s", err.Error()))
	}

	if len(s.Pass) > 0 {
		if _, err = s.Conn.Do("AUTH", s.Pass); err != nil {
			panic(fmt.Sprintf("redis auth error: %s", err.Error()))
		}
	}
}

func (s *MyRedis) autoReconnect() {

	for {

		if err := s.Conn.Err(); err != nil {
			fmt.Println(fmt.Sprintf("redis connection failure (%s)... attempting to restart: ", err))
			s.Conn.Close()
			s.connect()
		}
		time.Sleep(5000 * time.Millisecond)
	}
}
