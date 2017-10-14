package group

import (
	"github.com/boltdb/bolt"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/group/internal"
)

const (
	HYDRA_GROUP_BUCKET         = "GROUP"
	HYDRA_GROUP_MEMBERS_BUCKET = "GROUP_MEMBERS"
)

var (
	HYDRA_GROUP_BUCKET_BYTES         = []byte(HYDRA_GROUP_BUCKET)
	HYDRA_GROUP_MEMBERS_BUCKET_BYTES = []byte(HYDRA_GROUP_MEMBERS_BUCKET)
)

var (
	SOME_VALUE []byte = []byte{1}
)

type GroupManager struct {
	client *boltdbclient.Client
}

func NewGroupManager(client *boltdbclient.Client) (*GroupManager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_GROUP_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_GROUP_MEMBERS_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &GroupManager{
		client: client,
	}, nil
}

func (gm *GroupManager) createGroup(tx *bolt.Tx, g *group.Group) error {

	b := tx.Bucket(HYDRA_GROUP_BUCKET_BYTES)
	if v, err := internal.GroupMarshal(g); err != nil {
		return errors.WithStack(err)
	} else if err = b.Put([]byte(g.ID), v); err != nil {
		return errors.WithStack(err)
	}

	return gm.addGroupMembers(tx, g.ID, g.Members)
}

func (gm *GroupManager) CreateGroup(g *group.Group) error {
	return gm.client.WriteTransaction(func(tx *bolt.Tx) error {
		return gm.createGroup(tx, g)
	})
}

func (gm *GroupManager) getGroup(tx *bolt.Tx, id string) (*group.Group, error) {
	return gm.getGroupByte(tx, []byte(id))
}

func (gm *GroupManager) getGroupByte(tx *bolt.Tx, id []byte) (*group.Group, error) {
	b := tx.Bucket(HYDRA_GROUP_BUCKET_BYTES)

	if v := b.Get(id); v == nil {
		return nil, errors.Wrap(pkg.ErrNotFound, string(id))
	} else {
		var g group.Group
		if err := internal.GroupUnmarshal(v, &g); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return &g, nil
		}
	}
}

func (gm *GroupManager) GetGroup(id string) (*group.Group, error) {

	var g *group.Group

	return g, gm.client.ReadTransaction(func(tx *bolt.Tx) error {
		var err error
		g, err = gm.getGroup(tx, id)
		return err
	})
}

func (gm *GroupManager) DeleteGroup(id string) error {
	return gm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if g, err := gm.getGroup(tx, id); err != nil {
			return errors.WithStack(err)
		} else {
			b := tx.Bucket(HYDRA_GROUP_BUCKET_BYTES)
			if err := b.Delete([]byte(id)); err != nil {
				return errors.WithStack(err)
			}
			return gm.removeGroupMembers(tx, id, g.Members)
		}
		return nil
	})
}

func (gm *GroupManager) addGroupMembers(tx *bolt.Tx, group string, members []string) error {
	b := tx.Bucket(HYDRA_GROUP_MEMBERS_BUCKET_BYTES)
	for _, m := range members {
		if mb, err := b.CreateBucketIfNotExists([]byte(m)); err != nil {
			return errors.WithStack(err)
		} else {
			if err = mb.Put([]byte(group), SOME_VALUE); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

func (gm *GroupManager) AddGroupMembers(group string, members []string) error {
	return gm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if g, err := gm.getGroup(tx, group); err != nil {
			return errors.WithStack(err)
		} else {
			g.Members = append(g.Members, members...)
			if err = gm.createGroup(tx, g); err != nil {
				return errors.WithStack(err)
			}
			if err = gm.addGroupMembers(tx, group, members); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	})
}

func (gm *GroupManager) removeGroupMembers(tx *bolt.Tx, group string, members []string) error {
	b := tx.Bucket(HYDRA_GROUP_MEMBERS_BUCKET_BYTES)
	for _, m := range members {
		if mb := b.Bucket([]byte(m)); mb != nil {
			if err := mb.Delete([]byte(group)); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func (gm *GroupManager) RemoveGroupMembers(group string, members []string) error {
	return gm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if g, err := gm.getGroup(tx, group); err != nil {
			return errors.WithStack(err)
		} else {
			var subs []string
			for _, s := range g.Members {
				var remove bool
				for _, f := range members {
					if f == s {
						remove = true
						break
					}
				}
				if !remove {
					subs = append(subs, s)
				}
			}
			g.Members = subs
			if err = gm.createGroup(tx, g); err != nil {
				return errors.WithStack(err)
			}
			gm.removeGroupMembers(tx, group, members)
			return nil
		}
	})
}

func (gm *GroupManager) FindGroupsByMember(subject string) ([]group.Group, error) {
	var res = []group.Group{}
	return res, gm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_GROUP_MEMBERS_BUCKET_BYTES)
		if mb := b.Bucket([]byte(subject)); mb != nil {
			c := mb.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if g, err := gm.getGroupByte(tx, k); err != nil {
					return errors.WithStack(err)
				} else {
					res = append(res, *g)
				}
			}
		}
		return nil
	})
}

func (gm *GroupManager) FindGroupNames(member string) ([]string, error) {
	var res []string
	return res, gm.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(HYDRA_GROUP_MEMBERS_BUCKET_BYTES)
		if mb := b.Bucket([]byte(member)); mb != nil {
			c := mb.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				res = append(res, string(k))
			}
		}
		return nil
	})
}
