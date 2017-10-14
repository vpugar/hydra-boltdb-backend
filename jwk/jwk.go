package jwk

import (
	"github.com/boltdb/bolt"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/jwk/internal"
)

const (
	HYDRA_JWK_BUCKET = "JWK"
)

var (
	HYDRA_JWK_BUCKET_BYTES = []byte(HYDRA_JWK_BUCKET)
)

type JwkManager struct {
	client *boltdbclient.Client
}

func NewJwkManager(client *boltdbclient.Client) (*JwkManager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_JWK_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &JwkManager{
		client: client,
	}, nil
}

func (jm *JwkManager) addKey(bs *bolt.Bucket, set string, key *jose.JsonWebKey) error {
	kid := key.KeyID
	if b, err := bs.CreateBucketIfNotExists([]byte(kid)); err != nil {
		return errors.Wrapf(err, "set %v key %v", set, kid)
	} else if data, err := internal.JwkJsonWebKeyMarshal(key); err != nil {
		return errors.Wrapf(err, "set %v key %v", set, kid)
	} else {
		seq, _ := b.NextSequence()
		b.Put(boltdbclient.I2B(seq), data)
	}
	return nil
}

func (jm *JwkManager) AddKey(set string, key *jose.JsonWebKey) error {
	return jm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if bs, err := b.CreateBucketIfNotExists([]byte(set)); err != nil {
			return errors.Wrapf(err, "set %v", set)
		} else if err = jm.addKey(bs, set, key); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func (jm *JwkManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	return jm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if bs, err := b.CreateBucketIfNotExists([]byte(set)); err != nil {
			return errors.Wrapf(err, "set %v", set)
		} else {
			for _, key := range keys.Keys {
				if err = jm.addKey(bs, set, &key); err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		}
	})
}

func (jm *JwkManager) GetKey(set, kid string) (*jose.JsonWebKeySet, error) {
	var keySet *jose.JsonWebKeySet
	return keySet, jm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if bs := b.Bucket([]byte(set)); bs != nil {
			if bkey := bs.Bucket([]byte(kid)); bkey != nil {
				c := bkey.Cursor()
				var keys []jose.JsonWebKey
				for k, v := c.First(); k != nil; k, v = c.Next() {
					key := jose.JsonWebKey{}
					internal.JwkJsonWebKeyUnmarshal(v, &key)
					keys = append(keys, key)
				}
				keySet = &jose.JsonWebKeySet{
					Keys: keys,
				}
				return nil
			} else {
				return errors.Wrapf(pkg.ErrNotFound, "set %v key %v", set, kid)
			}
		} else {
			return errors.Wrapf(pkg.ErrNotFound, "set %v", set)
		}
	})
}

func (jm *JwkManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	var keySet *jose.JsonWebKeySet
	return keySet, jm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if bs := b.Bucket([]byte(set)); bs != nil {
			ckey := bs.Cursor()
			var keys []jose.JsonWebKey
			for k, _ := ckey.First(); k != nil; k, _ = ckey.Next() {
				c := bs.Bucket(k).Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					key := jose.JsonWebKey{}
					internal.JwkJsonWebKeyUnmarshal(v, &key)
					keys = append(keys, key)
				}
			}
			keySet = &jose.JsonWebKeySet{
				Keys: keys,
			}
		} else {
			return errors.Wrapf(pkg.ErrNotFound, "set %v", set)
		}
		return nil
	})
}

func (jm *JwkManager) DeleteKey(set, kid string) error {
	return jm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if bs := b.Bucket([]byte(set)); bs != nil {
			if err := bs.DeleteBucket([]byte(kid)); err != nil {
				return errors.Wrapf(err, "set %v key %v", set, kid)
			}
		} else {
			return errors.Wrapf(pkg.ErrNotFound, "set %v", set)
		}
		return nil
	})
}

func (jm *JwkManager) DeleteKeySet(set string) error {
	return jm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_JWK_BUCKET_BYTES)
		if err := b.DeleteBucket([]byte(set)); err != nil {
			return errors.Wrapf(err, "set %v", set)
		}
		return nil
	})
}
