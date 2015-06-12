# pezauth
authentication service for pez resources

[![wercker status](https://app.wercker.com/status/01d1f291e61f5edfc16f0b0ac182af8f/m/master "wercker status")](https://app.wercker.com/project/bykey/01d1f291e61f5edfc16f0b0ac182af8f)

[![GoDoc](https://godoc.org/github.com/pivotalservices/pezauth?status.png)](http://godoc.org/github.com/pivotalservices/pezauth)

## background
this is a combination of 2 products the auth service & the pez portal.
pez auth services is a central auth consumable that could be used to validate rest calls across all pez services
pez portal is the user facing web client (pez landing page) which will be a users point and click representation of pez

the two above services will be decoupled in the near future, but remain in this repo for now.

## How to run my pipeline
* requires the wercker cli
* requires a dockerhost (boot2docker)
```
$ git clone git@github.com:pivotalservices/pezauth.git
$ cd pezauth
$ ./runlocalbuild
```






