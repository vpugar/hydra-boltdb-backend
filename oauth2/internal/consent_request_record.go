package internal

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/ory/hydra/oauth2"
	"github.com/pkg/errors"
	"time"
)

//go:generate protoc --gofast_out=. consent_request_record.proto

func ConsentRequestMarshal(v *oauth2.ConsentRequest) ([]byte, error) {

	var accessTokenExtra []byte
	var idTokenExtra []byte

	if v.AccessTokenExtra != nil {
		if out, err := json.Marshal(v.AccessTokenExtra); err != nil {
			return nil, errors.WithStack(err)
		} else {
			accessTokenExtra = out
		}
	}

	if v.IDTokenExtra != nil {
		if out, err := json.Marshal(v.IDTokenExtra); err != nil {
			return nil, errors.WithStack(err)
		} else {
			idTokenExtra = out
		}
	}

	return proto.Marshal(&ConsentRequestRecord{
		ID:               v.ID,
		RequestedScopes:  v.RequestedScopes,
		ClientID:         v.ClientID,
		ExpiresAt:        v.ExpiresAt.UTC().UnixNano(),
		RedirectURL:      v.RedirectURL,
		CSRF:             v.CSRF,
		GrantedScopes:    v.GrantedScopes,
		Subject:          v.Subject,
		AccessTokenExtra: accessTokenExtra,
		IDTokenExtra:     idTokenExtra,
		Consent:          v.Consent,
		DenyReason:       v.DenyReason,
	})
}

func ConsentRequestUnmarshal(data []byte, v *oauth2.ConsentRequest) error {
	var pb ConsentRequestRecord
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	accessTokenExtra := new(map[string]interface{})
	if pb.AccessTokenExtra != nil {
		if err := json.Unmarshal(pb.AccessTokenExtra, accessTokenExtra); err != nil {
			return errors.WithStack(err)
		}
	}

	idTokenExtra := new(map[string]interface{})
	if pb.IDTokenExtra != nil {
		if err := json.Unmarshal(pb.IDTokenExtra, idTokenExtra); err != nil {
			return errors.WithStack(err)
		}
	}

	v.ID = pb.ID
	v.RequestedScopes = pb.RequestedScopes
	v.ClientID = pb.ClientID
	v.ExpiresAt = time.Unix(0, pb.ExpiresAt).UTC()
	v.RedirectURL = pb.RedirectURL
	v.CSRF = pb.CSRF
	v.GrantedScopes = pb.GrantedScopes
	v.Subject = pb.Subject
	v.AccessTokenExtra = *accessTokenExtra
	v.IDTokenExtra = *idTokenExtra
	v.Consent = pb.Consent
	v.DenyReason = pb.DenyReason

	return nil
}
