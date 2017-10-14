package internal_test

import (
	"github.com/ory/hydra/client"
	"github.com/vpugar/hydra-boltdb-backend/client/internal"
	"reflect"
	"testing"
)

func TestMarshalClient(t *testing.T) {
	v := client.Client{
		ID:                "ID",
		Name:              "Name",
		Secret:            "Secret",
		RedirectURIs:      []string{"1", "2"},
		GrantTypes:        []string{"3", "4"},
		ResponseTypes:     []string{"5", "6"},
		Scope:             "Scope",
		Owner:             "Owner",
		PolicyURI:         "PoliciURI",
		TermsOfServiceURI: "TermsOfServiceURI",
		ClientURI:         "ClientURI",
		LogoURI:           "LogoURI",
		Contacts:          []string{"7", "8"},
		Public:            true,
	}

	var other client.Client
	if buf, err := internal.ClientMarshal(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.ClientUnmarshal(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
