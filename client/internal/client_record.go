package internal

import (
	"github.com/gogo/protobuf/proto"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
)

//go:generate protoc --gofast_out=. client_record.proto

func ClientMarshal(v *client.Client) ([]byte, error) {
	return proto.Marshal(&ClientRecord{
		ID:                v.ID,
		Name:              v.Name,
		Secret:            v.Secret,
		RedirectURIs:      v.RedirectURIs,
		GrantTypes:        v.GrantTypes,
		ResponseTypes:     v.ResponseTypes,
		Scope:             v.Scope,
		Owner:             v.Owner,
		PolicyURI:         v.PolicyURI,
		TermsOfServiceURI: v.TermsOfServiceURI,
		ClientURI:         v.ClientURI,
		LogoURI:           v.LogoURI,
		Contacts:          v.Contacts,
		Public:            v.Public,
	})
}

func ClientUnmarshal(data []byte, v *client.Client) error {
	var pb ClientRecord
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	v.ID = pb.ID
	v.Name = pb.Name
	v.Secret = pb.Secret
	v.RedirectURIs = pb.RedirectURIs
	v.GrantTypes = pb.GrantTypes
	v.ResponseTypes = pb.ResponseTypes
	v.Scope = pb.Scope
	v.Owner = pb.Owner
	v.PolicyURI = pb.PolicyURI
	v.TermsOfServiceURI = pb.TermsOfServiceURI
	v.ClientURI = pb.ClientURI
	v.LogoURI = pb.LogoURI
	v.Contacts = pb.Contacts
	v.Public = pb.Public

	return nil
}

func FositeClientUnmarshal(data []byte, v *fosite.DefaultClient) error {
	var pb ClientRecord
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	v.ID = pb.ID
	v.Secret = []byte(pb.Secret)
	v.RedirectURIs = pb.RedirectURIs
	v.GrantTypes = pb.GrantTypes
	v.ResponseTypes = pb.ResponseTypes
	v.Public = pb.Public

	return nil
}
