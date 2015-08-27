# pezauth
authentication service for pez resources

[![wercker status](https://app.wercker.com/status/4ccfbbbb72ec786a0cc02dabc5de3f41/s/master "wercker status")](https://app.wercker.com/project/bykey/4ccfbbbb72ec786a0cc02dabc5de3f41)

[![GoDoc](https://godoc.org/github.com/pivotal-pez/pezauth?status.png)](http://godoc.org/github.com/pivotal-pez/pezauth)

## background
this is a combination of 2 products the auth service & the pez portal.
pez auth services is a central auth consumable that could be used to validate rest calls across all pez services
pez portal is the user facing web client (pez landing page) which will be a users point and click representation of pez

the two above services will be decoupled in the near future, but remain in this repo for now.



## Running tests / build pipeline locally

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

# make sure a docker host is running
$ boot2docker up && $(boot2docker shellinit)

# run the build pipeline locally, to test your code locally
$ ./testrunner

```


## Running locally for development

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

#lets bootstrap our repo as a local dev space
$ ./init_developer_environment

# make sure a docker host is running
$ boot2docker up && $(boot2docker shellinit)

# run the app locally using wercker magic
$ ./runlocaldeploy local_wercker_configs/myenv

$ echo "open ${DOCKER_HOST} in your browser to view this app locally"

```
