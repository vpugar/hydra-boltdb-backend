package backend

import (
	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	client2 "github.com/vpugar/hydra-boltdb-backend/client"
	fosite2 "github.com/vpugar/hydra-boltdb-backend/fosite"
	group2 "github.com/vpugar/hydra-boltdb-backend/group"
	jwk2 "github.com/vpugar/hydra-boltdb-backend/jwk"
	ladon2 "github.com/vpugar/hydra-boltdb-backend/ladon"
)

type Plugin struct {
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (c *Plugin) Connect(url string) (*sqlx.DB, error) {
	if db, err := sqlx.Open(driverName, url); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func (c *Plugin) NewClientManager(db *sqlx.DB, hasher fosite.Hasher) client.Manager {
	var cm client.Manager = &client2.ClientManager{}
	return cm
}

func (c *Plugin) NewGroupManager(db *sqlx.DB) group.Manager {
	var gm group.Manager = &group2.GroupManager{}
	return gm
}

func (c *Plugin) NewJWKManager(db *sqlx.DB, aead *jwk.AEAD) jwk.Manager {
	var jm jwk.Manager = &jwk2.JwkManager{}
	return jm
}

func (c *Plugin) NewOAuth2Manager(db *sqlx.DB, manager client.Manager, logger logrus.FieldLogger) pkg.FositeStorer {
	var om pkg.FositeStorer = &fosite2.Oauth2Manager{}
	return om
}

func (c *Plugin) NewPolicyManager(db *sqlx.DB) ladon.Manager {
	var lm ladon.Manager = &ladon2.LadonManager{}
	return lm
}
