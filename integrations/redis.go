package integrations

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
	"github.com/xchapter7x/lo"
)

//New - create a new redis integration wrapper
func (s *MyRedis) New(appEnv *cfenv.App) *MyRedis {
	redisName := os.Getenv("REDIS_SERVICE_NAME")
	redisHost := os.Getenv("REDIS_HOSTNAME_NAME")
	redisPass := os.Getenv("REDIS_PASSWORD_NAME")
	redisPort := os.Getenv("REDIS_PORT_NAME")
	redisService, err := appEnv.Services.WithName(redisName)
	s.Pass = redisService.Credentials[redisPass].(string)
	s.URI = fmt.Sprintf("%s:%d", redisService.Credentials[redisHost].(string), int(redisService.Credentials[redisPort].(float64)))

	if err != nil {
		panic(fmt.Sprintf("redis service name error: %s", err.Error()))
	}
	s.Pool = &redis.Pool{Dial: s.connect, MaxIdle: redisMaxIdle}
	return s
}

func (s *MyRedis) connect() (conn redis.Conn, err error) {

	if conn, err = redis.Dial("tcp", s.URI); err != nil {
		lo.G.Error(fmt.Sprintf("redis dial error: %s", err.Error()))

	} else {

		if len(s.Pass) > 0 {

			if _, err = conn.Do("AUTH", s.Pass); err != nil {
				lo.G.Error(fmt.Sprintf("redis auth error: %s", err.Error()))
			}
		}
	}
	return
}
