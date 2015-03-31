package main

import (
	"github.com/go-martini/martini"
	pez "github.com/pivotalservices/pezauth/service"
)

func main() {
	m := martini.Classic()
	pez.InitRoutes(m)
	m.Run()
}
