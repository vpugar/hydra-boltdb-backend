package ladon

import (
	"github.com/boltdb/bolt"
	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/ladon/internal"
)

const (
	HYDRA_LADOM_BUCKET = "LADOM"
)

var (
	HYDRA_LADOM_BUCKET_BYTES = []byte(HYDRA_LADOM_BUCKET)
)

type LadonManager struct {
	client *boltdbclient.Client
}

func NewLadonManager(client *boltdbclient.Client) (*LadonManager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_LADOM_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &LadonManager{
		client: client,
	}, nil
}

// Create persists the policy.
func (lm *LadonManager) Create(policy ladon.Policy) error {
	return lm.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_LADOM_BUCKET_BYTES)
		if v, err := internal.PolicyMarshal(&policy); err != nil {
			return errors.WithStack(err)
		} else if err := b.Put([]byte(policy.GetID()), v); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update updates an existing policy.
func (lm *LadonManager) Update(policy ladon.Policy) error {
	return lm.Create(policy)
}

func (lm *LadonManager) get(tx *bolt.Tx, id string) (ladon.Policy, error) {
	b := tx.Bucket(HYDRA_LADOM_BUCKET_BYTES)

	if v := b.Get([]byte(id)); v == nil {
		return nil, errors.Wrap(pkg.ErrNotFound, id)
	} else {
		var p ladon.DefaultPolicy
		if err := internal.PolicyUnmarshal(v, &p); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return &p, nil
		}
	}
}

// Get retrieves a policy.
func (lm *LadonManager) Get(id string) (ladon.Policy, error) {
	var p ladon.Policy
	return p, lm.client.ReadTransaction(func(tx *bolt.Tx) error {
		var err error
		p, err = lm.get(tx, id)
		return err
	})
}

// Delete removes a policy.
func (lm *LadonManager) Delete(id string) error {
	return lm.client.DeleteWithTransaction(HYDRA_LADOM_BUCKET_BYTES, id)
}

// GetAll retrieves all policies.
func (lm *LadonManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	var policies ladon.Policies
	return policies, lm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_LADOM_BUCKET_BYTES)
		c := b.Cursor()
		var count int64
		max := limit + offset - 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var instance ladon.DefaultPolicy
			if count >= offset {
				if err := internal.PolicyUnmarshal(v, &instance); err != nil {
					return errors.WithStack(err)
				} else {
					policies = append(policies, &instance)
				}
				if max <= count {
					break
				}
			}
			count++
		}
		return nil
	})
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
// FIXME check in manager_sql how to filter out subjects
func (lm *LadonManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	var policies ladon.Policies
	return policies, lm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_LADOM_BUCKET_BYTES)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var instance ladon.DefaultPolicy
			if err := internal.PolicyUnmarshal(v, &instance); err != nil {
				return errors.WithStack(err)
			} else {
				policies = append(policies, &instance)
			}
		}
		return nil
	})
}
