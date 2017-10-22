package internal_test

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/vpugar/hydra-boltdb-backend/fosite/internal"
)

var clientManager client.Manager = &client.MemoryManager{
	Clients: map[string]client.Client{"foobar": {ID: "foobar"}},
	Hasher:  &fosite.BCrypt{},
}

func TestMarshalOauth2Record(t *testing.T) {
	signature := "SIG1"

	form := url.Values{
		"k1": []string{"v1"},
		"k2": []string{"v2"},
	}

	v := fosite.Request{
		ID:          "ID",
		RequestedAt: time.Now().UTC(),
		Client: &fosite.DefaultClient{
			ID: "foobar",
		},
		Scopes:        []string{"3", "4"},
		GrantedScopes: []string{"5", "6"},
		Form:          form,
		Session: &fosite.DefaultSession{
			Username: "user1",
			Subject:  "sub1",
		},
	}

	var other fosite.Request
	if buf, err := internal.OAuth2RequesterMarshal(signature, &v); err != nil {
		t.Fatal(err)
	} else if err := internal.OAuth2RequesterUnmarshal(buf, &other, &fosite.DefaultSession{}, clientManager); err != nil {
		t.Fatal(err)
	} else {
		if !reflect.DeepEqual(v.ID, other.ID) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
		if !reflect.DeepEqual(v.RequestedAt, other.RequestedAt) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
		if !reflect.DeepEqual(v.Scopes, other.Scopes) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
		if !reflect.DeepEqual(v.GrantedScopes, other.GrantedScopes) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
		if !reflect.DeepEqual(v.Form, other.Form) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
		if !reflect.DeepEqual(v.Session, other.Session) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other.Session, v.Session)
		}
		if !reflect.DeepEqual(v.Client.GetID(), other.Client.GetID()) {
			t.Fatalf("unexpected copy:\n\t%#v\n\t%#v", other, v)
		}
	}
}
