package client

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/imdario/mergo"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/client/internal"
)

const (
	HYDRA_CLIENT_BUCKET = "CLIENT"
)

var (
	HYDRA_CLIENT_BUCKET_BYTES = []byte(HYDRA_CLIENT_BUCKET)
)

type ClientManager struct {
	client *boltdbclient.Client
	hasher fosite.Hasher
}

func NewClientManager(client *boltdbclient.Client, hasher fosite.Hasher) (*ClientManager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_CLIENT_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ClientManager{
		client: client,
		hasher: hasher,
	}, nil
}

func (cm *ClientManager) Authenticate(id string, secret []byte) (*client.Client, error) {
	var c *client.Client
	return c, cm.client.ReadTransaction(func(tx *bolt.Tx) error {
		var err error
		if c, err = cm.getClient(tx, id); err != nil {
			return errors.WithStack(err)
		} else if err := cm.hasher.Compare(c.GetHashedSecret(), secret); err != nil {
			return errors.WithStack(err)
		} else {
			return nil
		}
	})
}

func (cm *ClientManager) getClient(tx *bolt.Tx, id string) (*client.Client, error) {
	b := tx.Bucket(HYDRA_CLIENT_BUCKET_BYTES)

	if v := b.Get([]byte(id)); v == nil {
		return nil, errors.Wrap(pkg.ErrNotFound, id)
	} else {
		var c client.Client
		if err := internal.ClientUnmarshal(v, &c); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return &c, nil
		}
	}
}

func (cm *ClientManager) GetConcreteClient(id string) (*client.Client, error) {
	var c *client.Client
	return c, cm.client.ReadTransaction(func(tx *bolt.Tx) error {
		var err error
		c, err = cm.getClient(tx, id)
		return err
	})
}

func (cm *ClientManager) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return cm.GetConcreteClient(id)
}

func (cm *ClientManager) CreateClient(c *client.Client) error {

	if c.ID == "" {
		c.ID = uuid.New()
	}

	hash, err := cm.hasher.Hash([]byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(hash)

	return cm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_CLIENT_BUCKET_BYTES)

		if v, err := internal.ClientMarshal(c); err != nil {
			return errors.WithStack(err)
		} else if err := b.Put([]byte(c.ID), v); err != nil {
			return errors.WithStack(err)
		} else {
			return nil
		}
	})
}

func (cm *ClientManager) UpdateClient(c *client.Client) error {
	return cm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if o, err := cm.getClient(tx, c.ID); err != nil {
			return errors.WithStack(err)
		} else {
			if c.Secret == "" {
				c.Secret = string(o.GetHashedSecret())
			} else {
				h, err := cm.hasher.Hash([]byte(c.Secret))
				if err != nil {
					return errors.WithStack(err)
				}
				c.Secret = string(h)
			}

			if err := mergo.Merge(c, o); err != nil {
				return errors.WithStack(err)
			}

			b := tx.Bucket(HYDRA_CLIENT_BUCKET_BYTES)

			if v, err := internal.ClientMarshal(c); err != nil {
				return errors.WithStack(err)
			} else if err := b.Put([]byte(c.ID), v); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

func (cm *ClientManager) DeleteClient(id string) error {
	return cm.client.DeleteWithTransaction(HYDRA_CLIENT_BUCKET_BYTES, id)
}

func (cm *ClientManager) GetClients() (map[string]client.Client, error) {

	clients := make(map[string]client.Client)

	err := cm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_CLIENT_BUCKET_BYTES)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var clientInstance client.Client
			if err := internal.ClientUnmarshal(v, &clientInstance); err != nil {
				return errors.WithStack(err)
			} else {
				clients[clientInstance.ID] = clientInstance
			}
		}

		return nil
	})

	return clients, err
}
