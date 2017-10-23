# Hydra BoltDB Backend

<p align="left">
    <a href="https://travis-ci.org/vpugar/hydra-boltdb-backend"><img src="https://travis-ci.org/vpugar/hydra-boltdb-backend.svg?branch=master" alt="Build Status"></a>
    <a href="https://coveralls.io/github/vpugar/hydra-boltdb-backend?branch=master"><img src="https://coveralls.io/repos/vpugar/hydra-boltdb-backend/badge.svg?branch=master&service=github" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/vpugar/hydra-boltdb-backend"><img src="https://goreportcard.com/badge/github.com/vpugar/hydra-boltdb-backend" alt="Go Report Card"></a>
</p>

under construction :)

## Building with hydra

To integrate with hydra, it is used hydra version from https://github.com/vpugar/hydra.
The version has additional factory parts to use boltdb implementation.

### Getting hydra

Create missing directories if needed.

    cd $GOPATH/src/github.com/ory/
    git clone https://github.com/vpugar/hydra.git

### Building hydra

    cd $GOPATH/src/github.com/ory/hydra
    glide install
    go install github.com/ory/hydra

## Startup hydra (in development mode) with boltdb

	export DATABASE_URL=boltdb://hydra-boltdb.db
	export SYSTEM_SECRET=zF8JIbGE4uBBOkNHDpxH5VUCmmMqkQ6V
	export CONSENT_URL=http://localhost:3000/login
	cd $GOPATH/bin;./hydra host --dangerous-force-http

## Using with consent app hydra-idp-react

### Configuration

#### hydra connect

Run

    hydra connect

Check hydra log

    INFO[0002] client_id: c940dbd8-9ffd-4f6e-ae99-cc7536f06238
    INFO[0002] client_secret: Wx5NL_yB7qGFTy_-

Example:
* cluster_url: http://localhost:4444
* client_id: c940dbd8-9ffd-4f6e-ae99-cc7536f06238
* client_secret: Wx5NL_yB7qGFTy_-

Test running hydra with

	hydra token client --skip-tls-verify
	hydra token validate --skip-tls-verify $(hydra token client --skip-tls-verify)

Move created $HOME/.hydra.yml to $HOME/.custom_hydra.yml

#### consent app client setup

    hydra clients create --skip-tls-verify \
      --id consent-app \
      --secret consent-secret \
      --name "Consent App Client" \
      --grant-types client_credentials \
      --response-types token \
      --allowed-scopes hydra.consent \
      --config $HOME/.custom_hydra.yml

    hydra policies create --skip-tls-verify \
      --actions get,accept,reject \
      --description "Allow consent-app to manage OAuth2 consent requests." \
      --allow \
      --id consent-app-policy \
      --resources "rn:hydra:oauth2:consent:requests:<.*>" \
      --subjects consent-app \
      --config $HOME/.custom_hydra.yml

#### consent app build

Use
* hydra-idp-react (https://github.com/ory/hydra-idp-react) - NOT WORKING - currently works with challenge not consent parameter
* hydra-consent-app-express (https://github.com/ory/hydra-consent-app-express) - WORKING OK

Example for build:

    cd $SOMEWHERE
	git clone https://github.com/ory/hydra-idp-react
	npm i

#### consent app run

##### hydra-idp-react

Note: remove $HOME/.hydra.yml!!!

	cd $SOMEWHERE/hydra-idp-react/
	export HYDRA_CLIENT_ID=consent-app
	export HYDRA_CLIENT_SECRET=consent-secret
	export HYDRA_URL=http://localhost:4444
	export NODE_TLS_REJECT_UNAUTHORIZED=0
	npm run dev

Login URL: http://localhost:3000/

##### hydra-consent-app-express

Note: remove $HOME/.hydra.yml!!!

	cd $SOMEWHERE/hydra-consent-app-express/
	export HYDRA_CLIENT_ID=consent-app
	export HYDRA_CLIENT_SECRET=consent-secret
	export HYDRA_URL=http://localhost:4444
	export NODE_TLS_REJECT_UNAUTHORIZED=0
	npm start

Login URL: http://localhost:3000/

#### app client setup

    hydra clients create --skip-tls-verify \
      --id some-consumer \
      --secret consumer-secret \
      --grant-types authorization_code,refresh_token,client_credentials,implicit \
      --response-types token,code,id_token \
      --allowed-scopes openid,offline,hydra.clients \
      --callbacks http://localhost:4445/callback \
      --config $HOME/.custom_hydra.yml

    hydra policies create --skip-tls-verify \
      --actions get \
      --description "Allow everyone to read the OpenID Connect ID Token public key" \
      --allow \
      --id openid-id_token-policy \
      --resources rn:hydra:keys:hydra.openid.id-token:public \
      --subjects "<.*>" \
      --config $HOME/.custom_hydra.yml


#### login flow

    hydra token user --skip-tls-verify \
      --auth-url http://localhost:4444/oauth2/auth \
      --token-url http://localhost:4444/oauth2/token \
      --id some-consumer \
      --secret consumer-secret \
      --scopes openid,offline,hydra.clients \
      --redirect http://localhost:4445/callback

## NOTE

Initially developed as plugin with, lets say "workarround" to fake SQL connection with boltdb connection (implementation in plugin package).
But couldn't load plugin because of interfaces that defined in hydra from vendored dependencies.
Interfaces in that case have not same signature as one in the plugin.

More details is following issue: https://github.com/golang/go/issues/20481.
