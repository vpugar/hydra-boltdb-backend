package oauth2

import (
	"github.com/boltdb/bolt"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/oauth2/internal"
)

const (
	HYDRA_CONSENT_BUCKET = "CONSENT"
)

var (
	HYDRA_CONSENT_BUCKET_BYTES = []byte(HYDRA_CONSENT_BUCKET)
)

type ConsentRequestManager struct {
	client *boltdbclient.Client
}

func NewConsentRequestManager(client *boltdbclient.Client) (*ConsentRequestManager, error) {
	// Initialize top-level buckets.
	if err := client.InitEntity(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(HYDRA_CONSENT_BUCKET_BYTES); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ConsentRequestManager{
		client: client,
	}, nil
}

func (crm *ConsentRequestManager) createConsentRequest(tx *bolt.Tx, request *oauth2.ConsentRequest) error {
	b := tx.Bucket(HYDRA_CONSENT_BUCKET_BYTES)
	if v, err := internal.ConsentRequestMarshal(request); err != nil {
		return errors.WithStack(err)
	} else if err = b.Put([]byte(request.ID), v); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (crm *ConsentRequestManager) PersistConsentRequest(request *oauth2.ConsentRequest) error {
	return crm.client.WriteTransaction(func(tx *bolt.Tx) error {
		return crm.createConsentRequest(tx, request)
	})
}

func (crm *ConsentRequestManager) AcceptConsentRequest(id string, payload *oauth2.AcceptConsentRequestPayload) error {
	return crm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if r, err := crm.get(tx, id); err != nil {
			return errors.WithStack(err)
		} else {
			r.Subject = payload.Subject
			r.AccessTokenExtra = payload.AccessTokenExtra
			r.IDTokenExtra = payload.IDTokenExtra
			r.Consent = oauth2.ConsentRequestAccepted
			r.GrantedScopes = payload.GrantScopes
			return crm.createConsentRequest(tx, r)
		}
	})
}

func (crm *ConsentRequestManager) RejectConsentRequest(id string, payload *oauth2.RejectConsentRequestPayload) error {
	return crm.client.WriteTransaction(func(tx *bolt.Tx) error {
		if r, err := crm.get(tx, id); err != nil {
			return errors.WithStack(err)
		} else {
			r.Consent = oauth2.ConsentRequestRejected
			r.DenyReason = payload.Reason
			return crm.createConsentRequest(tx, r)
		}
	})
}

func (crm *ConsentRequestManager) get(tx *bolt.Tx, id string) (*oauth2.ConsentRequest, error) {
	b := tx.Bucket(HYDRA_CONSENT_BUCKET_BYTES)

	if v := b.Get([]byte(id)); v == nil {
		return nil, errors.Wrap(pkg.ErrNotFound, id)
	} else {
		var p oauth2.ConsentRequest
		if err := internal.ConsentRequestUnmarshal(v, &p); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return &p, nil
		}
	}
}

func (crm *ConsentRequestManager) GetConsentRequest(id string) (*oauth2.ConsentRequest, error) {
	var cr *oauth2.ConsentRequest
	return cr, crm.client.ReadTransaction(func(tx *bolt.Tx) error {
		var err error
		cr, err = crm.get(tx, id)
		return err
	})
}
