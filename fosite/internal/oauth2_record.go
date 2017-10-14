package internal

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

//go:generate protoc --gofast_out=. oauth2_record.proto

func OAuth2RequesterMarshal(signature string, v fosite.Requester) ([]byte, error) {

	session, err := json.Marshal(v.GetSession())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	form := map[string]*Val{}
	for k, v := range v.GetRequestForm() {
		form[k] = &Val{v}
	}

	return proto.Marshal(&Oauth2Record{
		Signature:     signature,
		Request:       v.GetID(),
		RequestedAt:   v.GetRequestedAt().UTC().UnixNano(),
		Client:        v.GetClient().GetID(),
		Scopes:        v.GetRequestedScopes(),
		GrantedScopes: v.GetGrantedScopes(),
		Form:          form,
		Session:       session,
	})
}

func OAuth2RequesterUnmarshal(data []byte, v *fosite.Request, session fosite.Session, cm fosite.ClientManager) error {
	var pb Oauth2Record

	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	if session != nil {
		if err := json.Unmarshal(pb.Session, session); err != nil {
			return errors.WithStack(err)
		}
	}

	if cm != nil {
		if c, err := cm.GetClient(context.Background(), pb.Client); err != nil {
			return errors.WithStack(err)
		} else {
			v.Client = c
		}
	}

	form := url.Values{}
	for k, v := range pb.Form {
		form[k] = v.Val
	}

	v.ID = pb.Request
	v.RequestedAt = time.Unix(0, pb.RequestedAt).UTC()
	v.Scopes = pb.Scopes
	v.GrantedScopes = pb.GrantedScopes
	v.Form = form
	v.Session = session

	return nil
}
