package internal_test

import (
	"github.com/ory/hydra/oauth2"
	"github.com/vpugar/hydra-boltdb-backend/oauth2/internal"
	"reflect"
	"testing"
	"time"
)

func TestMarshalConsentRequest(t *testing.T) {

	accessTokenExtra := map[string]interface{}{"11": "11", "12": "12"}
	idTokenExtra := map[string]interface{}{"11": "11", "12": "12"}

	v := oauth2.ConsentRequest{
		ID:               "ID",
		RequestedScopes:  []string{"11", "12"},
		ClientID:         "ClientID",
		ExpiresAt:        time.Now().UTC(),
		RedirectURL:      "RedirectURL",
		CSRF:             "CSRF",
		GrantedScopes:    []string{"21", "22"},
		Subject:          "Subject",
		AccessTokenExtra: accessTokenExtra,
		IDTokenExtra:     idTokenExtra,
		Consent:          "Consent",
		DenyReason:       "DenyReason",
	}

	var other oauth2.ConsentRequest
	if buf, err := internal.ConsentRequestMarshal(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.ConsentRequestUnmarshal(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: \n %#v\n %#v", other, v)
	}
}

func TestMarshalConsentRequestWithNullMaps(t *testing.T) {

	v := oauth2.ConsentRequest{
		ID:               "ID",
		RequestedScopes:  []string{"11", "12"},
		ClientID:         "ClientID",
		ExpiresAt:        time.Now().UTC(),
		RedirectURL:      "RedirectURL",
		CSRF:             "CSRF",
		GrantedScopes:    []string{"21", "22"},
		Subject:          "Subject",
		AccessTokenExtra: nil,
		IDTokenExtra:     nil,
		Consent:          "Consent",
		DenyReason:       "DenyReason",
	}

	var other oauth2.ConsentRequest
	if buf, err := internal.ConsentRequestMarshal(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.ConsentRequestUnmarshal(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: \n %#v\n %#v", other, v)
	}
}
