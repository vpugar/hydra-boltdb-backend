package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	client2 "github.com/vpugar/hydra-boltdb-backend/client"
	fosite2 "github.com/vpugar/hydra-boltdb-backend/fosite"
	group2 "github.com/vpugar/hydra-boltdb-backend/group"
	jwk2 "github.com/vpugar/hydra-boltdb-backend/jwk"
	ladon2 "github.com/vpugar/hydra-boltdb-backend/ladon"
	"github.com/vpugar/hydra-boltdb-backend/plugin/internal"
)

func Connect(url string) (*sqlx.DB, error) {
	if db, err := sqlx.Open(internal.DriverName, url); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func NewClientManager(db *sqlx.DB, hasher fosite.Hasher) client.Manager {
	var cm client.Manager = &client2.ClientManager{}
	return cm
}

func NewGroupManager(db *sqlx.DB) group.Manager {
	var gm group.Manager = &group2.GroupManager{}
	return gm
}

func NewJWKManager(db *sqlx.DB, aead *jwk.AEAD) jwk.Manager {
	var jm jwk.Manager = &jwk2.JwkManager{}
	return jm
}

func NewOAuth2Manager(db *sqlx.DB, manager client.Manager) pkg.FositeStorer {
	var om pkg.FositeStorer = &fosite2.Oauth2Manager{}
	return om
}

func NewPolicyManager(db *sqlx.DB) ladon.Manager {
	var lm ladon.Manager = &ladon2.LadonManager{}
	return lm
}
