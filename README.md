# pezauth
authentication service for pez resources

[![wercker status](https://app.wercker.com/status/4ccfbbbb72ec786a0cc02dabc5de3f41/s/master "wercker status")](https://app.wercker.com/project/bykey/4ccfbbbb72ec786a0cc02dabc5de3f41)

[![GoDoc](https://godoc.org/github.com/pivotal-pez/pezauth?status.png)](http://godoc.org/github.com/pivotal-pez/pezauth)

## background
this is a combination of 2 products the auth service & the pez portal.
pez auth services is a central auth consumable that could be used to validate rest calls across all pez services
pez portal is the user facing web client (pez landing page) which will be a users point and click representation of pez

the two above services will be decoupled in the near future, but remain in this repo for now.

## [How to run local deployment](docs\RUN_PEZAUTH_LOCAL.md)
