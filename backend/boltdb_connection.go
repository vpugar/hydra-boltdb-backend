package backend

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	oauth2_2 "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/vpugar/boltdbclient"
	client2 "github.com/vpugar/hydra-boltdb-backend/client"
	fosite2 "github.com/vpugar/hydra-boltdb-backend/fosite"
	group2 "github.com/vpugar/hydra-boltdb-backend/group"
	jwk2 "github.com/vpugar/hydra-boltdb-backend/jwk"
	ladon2 "github.com/vpugar/hydra-boltdb-backend/ladon"
	"github.com/vpugar/hydra-boltdb-backend/oauth2"
)

type BoltdbConnection struct {
	Client *boltdbclient.Client
}

func NewBoltdbConnection(databaseURL string) *BoltdbConnection {
	c := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: databaseURL,
	})
	return &BoltdbConnection{Client: c}
}

func (c *BoltdbConnection) Connect() error {
	if _, err := c.Client.Open(); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *BoltdbConnection) Disconnect() error {
	return c.Client.Close()
}

func (c *BoltdbConnection) NewClientManager(hasher fosite.Hasher) (client.Manager, error) {
	if cm, err := client2.NewClientManager(c.Client, hasher); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (c *BoltdbConnection) NewGroupManager() (group.Manager, error) {
	if cm, err := group2.NewGroupManager(c.Client); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (c *BoltdbConnection) NewJWKManager(aead *jwk.AEAD) (jwk.Manager, error) {
	// TODO aead *jwk.AEAD encryption of the storage
	if cm, err := jwk2.NewJwkManager(c.Client); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (c *BoltdbConnection) NewOAuth2Manager(manager client.Manager) (pkg.FositeStorer, error) {
	if cm, err := fosite2.NewOauth2Manager(c.Client, manager); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (c *BoltdbConnection) NewPolicyManager() (ladon.Manager, error) {
	if cm, err := ladon2.NewLadonManager(c.Client); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}

func (c *BoltdbConnection) NewConsentRequestManager() (oauth2_2.ConsentRequestManager, error) {
	if cm, err := oauth2.NewConsentRequestManager(c.Client); err != nil {
		return nil, err
	} else {
		return cm, nil
	}
}
