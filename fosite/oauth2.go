package fosite

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
	//"github.com/vpugar/hydra-boltdb-backend/client"
	client2 "github.com/ory/hydra/client"
	"github.com/vpugar/hydra-boltdb-backend/fosite/internal"
)

const (
	HYDRA_OAUTH2_OIDC_BUCKET       = "OAUTH2_OIDC"
	HYDRA_OAUTH2_ACCESS_BUCKET     = "OAUTH2_ACCESS"
	HYDRA_OAUTH2_ACCESS_ID_BUCKET  = "OAUTH2_ACCESS_ID"
	HYDRA_OAUTH2_REFRESH_BUCKET    = "OAUTH2_REFRESH"
	HYDRA_OAUTH2_REFRESH_ID_BUCKET = "OAUTH2_REFRESH_ID"
	HYDRA_OAUTH2_CODE_BUCKET       = "OAUTH2_CODE"
)

var (
	EMPTY_BYTE = []byte{1}
)

var (
	HYDRA_OAUTH2_OIDC_BUCKET_BYTES       = []byte(HYDRA_OAUTH2_OIDC_BUCKET)
	HYDRA_OAUTH2_ACCESS_BUCKET_BYTES     = []byte(HYDRA_OAUTH2_ACCESS_BUCKET)
	HYDRA_OAUTH2_ACCESS_ID_BUCKET_BYTES  = []byte(HYDRA_OAUTH2_ACCESS_ID_BUCKET)
	HYDRA_OAUTH2_REFRESH_BUCKET_BYTES    = []byte(HYDRA_OAUTH2_REFRESH_BUCKET)
	HYDRA_OAUTH2_REFRESH_ID_BUCKET_BYTES = []byte(HYDRA_OAUTH2_REFRESH_ID_BUCKET)
	HYDRA_OAUTH2_CODE_BUCKET_BYTES       = []byte(HYDRA_OAUTH2_CODE_BUCKET)
)

type Oauth2Manager struct {
	client        *boltdbclient.Client
	clientManager client2.Manager
}

func NewOauth2Manager(client *boltdbclient.Client, clientManager client2.Manager) (*Oauth2Manager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_OIDC_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_ACCESS_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_ACCESS_ID_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_REFRESH_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_REFRESH_ID_BUCKET_BYTES); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(HYDRA_OAUTH2_CODE_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Oauth2Manager{
		client:        client,
		clientManager: clientManager,
	}, nil
}

func (om *Oauth2Manager) createSession(signature string, requester fosite.Requester, bucketName []byte, bucketByIdName []byte) error {
	return om.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if val, err := internal.OAuth2RequesterMarshal(signature, requester); err != nil {
			return errors.WithStack(err)
		} else {
			b.Put([]byte(signature), val)
		}
		if bucketByIdName != nil {
			bByIds := tx.Bucket(bucketByIdName)
			if bId, err := bByIds.CreateBucketIfNotExists([]byte(requester.GetID())); err != nil {
				return errors.WithStack(err)
			} else {
				bId.Put([]byte(signature), EMPTY_BYTE)
			}
		}
		return nil
	})
}

func (om *Oauth2Manager) findSessionBySignature(signature string, session fosite.Session, bucketName []byte) (fosite.Requester, error) {
	r := &fosite.Request{}
	return r, om.client.ReadTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if data := b.Get([]byte(signature)); data == nil {
			return errors.Wrap(pkg.ErrNotFound, signature)
		} else {
			return internal.OAuth2RequesterUnmarshal(data, r, session, om.clientManager)
		}
	})
}

func (om *Oauth2Manager) deleteSession(signature string, bucketName []byte, bucketByIdName []byte) error {
	return om.client.WriteTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		r := b.Get([]byte(signature))
		if r != nil {
			var req fosite.Request
			if err := internal.OAuth2RequesterUnmarshal(r, &req, nil, nil); err != nil {
				return errors.WithStack(err)
			}
			if err := b.Delete([]byte(signature)); err != nil {
				return errors.WithStack(err)
			}
			if bucketByIdName != nil {
				bByIds := tx.Bucket(bucketByIdName)
				if err := bByIds.DeleteBucket([]byte(req.ID)); err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		} else {
			return errors.Wrap(pkg.ErrNotFound, signature)
		}
	})
}

func (om *Oauth2Manager) revokeSession(id string, bucketName []byte, bucketByIdName []byte) error {
	return om.client.WriteTransaction(func(tx *bolt.Tx) error {
		bByIds := tx.Bucket(bucketByIdName)

		if bId := bByIds.Bucket([]byte(id)); bId == nil {
			return errors.Wrap(pkg.ErrNotFound, id)
		} else {
			bSignature := tx.Bucket(bucketName)
			c := bId.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if err := bSignature.Delete(k); err != nil {
					return errors.WithStack(err)
				}
			}
			if err := bByIds.DeleteBucket([]byte(id)); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	})
}

func (om *Oauth2Manager) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return om.createSession(signature, request, HYDRA_OAUTH2_ACCESS_BUCKET_BYTES, HYDRA_OAUTH2_ACCESS_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) CreateImplicitAccessTokenSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return om.CreateAccessTokenSession(ctx, signature, requester)
}

func (om *Oauth2Manager) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return om.findSessionBySignature(signature, session, HYDRA_OAUTH2_ACCESS_BUCKET_BYTES)
}

func (om *Oauth2Manager) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return om.deleteSession(signature, HYDRA_OAUTH2_ACCESS_BUCKET_BYTES, HYDRA_OAUTH2_ACCESS_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	return om.createSession(code, request, HYDRA_OAUTH2_CODE_BUCKET_BYTES, nil)
}

func (om *Oauth2Manager) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	return om.findSessionBySignature(code, session, HYDRA_OAUTH2_CODE_BUCKET_BYTES)
}

func (om *Oauth2Manager) DeleteAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	return om.deleteSession(code, HYDRA_OAUTH2_CODE_BUCKET_BYTES, nil)
}

func (om *Oauth2Manager) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := om.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := om.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := om.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (om *Oauth2Manager) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return om.createSession(signature, request, HYDRA_OAUTH2_REFRESH_BUCKET_BYTES, HYDRA_OAUTH2_REFRESH_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return om.findSessionBySignature(signature, session, HYDRA_OAUTH2_REFRESH_BUCKET_BYTES)
}

func (om *Oauth2Manager) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return om.deleteSession(signature, HYDRA_OAUTH2_REFRESH_BUCKET_BYTES, HYDRA_OAUTH2_REFRESH_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := om.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err := om.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := om.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (om *Oauth2Manager) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {
	return om.createSession(authorizeCode, requester, HYDRA_OAUTH2_OIDC_BUCKET_BYTES, nil)
}

func (om *Oauth2Manager) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	return om.findSessionBySignature(authorizeCode, requester.GetSession(), HYDRA_OAUTH2_OIDC_BUCKET_BYTES)
}

func (om *Oauth2Manager) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	return om.deleteSession(authorizeCode, HYDRA_OAUTH2_OIDC_BUCKET_BYTES, nil)
}

func (om *Oauth2Manager) RevokeRefreshToken(ctx context.Context, requestID string) error {
	return om.revokeSession(requestID, HYDRA_OAUTH2_REFRESH_BUCKET_BYTES, HYDRA_OAUTH2_REFRESH_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) RevokeAccessToken(ctx context.Context, requestID string) error {
	return om.revokeSession(requestID, HYDRA_OAUTH2_ACCESS_BUCKET_BYTES, HYDRA_OAUTH2_ACCESS_ID_BUCKET_BYTES)
}

func (om *Oauth2Manager) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return om.clientManager.GetClient(ctx, id)
}
