package oauth2_test

import (
	"flag"
	"fmt"
	oauth2_2 "github.com/ory/hydra/oauth2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/oauth2"
	"os"
	"testing"
	"time"
)

const (
	name   = "bolt"
	dbName = "oauth2.test.db"
)

var consentRequestManager oauth2_2.ConsentRequestManager

func startTest(m *testing.M) int {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	if crm, err := oauth2.NewConsentRequestManager(boltClient); err != nil {
		fmt.Println("FAIL")
		return 1
	} else {
		consentRequestManager = crm
	}
	return m.Run()
}

func stopTest(boltClient *boltdbclient.Client) {
	defer os.Remove(dbName)
	if err := boltClient.Close(); err != nil {
		fmt.Println("FAIL", err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(startTest(m))
}

func TestInterface(t *testing.T) {
	var _ oauth2_2.ConsentRequestManager = &oauth2.ConsentRequestManager{}
}

func TestConsentRequestManagerReadWrite(t *testing.T) {
	req := &oauth2_2.ConsentRequest{
		ID:               "id-1",
		ClientID:         "client-id",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Minute),
		Consent:          oauth2_2.ConsentRequestAccepted,
		DenyReason:       "some reason",
		AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
		IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
		RedirectURL:      "https://redirect-me/foo",
		Subject:          "Peter",
	}

	k := name
	m := consentRequestManager
	t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
		_, err := m.GetConsentRequest("1234")
		assert.Error(t, err)

		require.NoError(t, m.PersistConsentRequest(req))

		got, err := m.GetConsentRequest(req.ID)
		require.NoError(t, err)

		require.Equal(t, req.ExpiresAt.Unix(), got.ExpiresAt.Unix())
		got.ExpiresAt = req.ExpiresAt
		assert.EqualValues(t, req, got)
	})
}

func TestConsentRequestManagerUpdate(t *testing.T) {
	req := &oauth2_2.ConsentRequest{
		ID:               "id-2",
		ClientID:         "client-id",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Minute),
		Consent:          oauth2_2.ConsentRequestRejected,
		DenyReason:       "some reason",
		AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
		IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
		RedirectURL:      "https://redirect-me/foo",
		Subject:          "Peter",
	}

	k := name
	m := consentRequestManager
	t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
		require.NoError(t, m.PersistConsentRequest(req))

		got, err := m.GetConsentRequest(req.ID)
		require.NoError(t, err)
		assert.False(t, got.IsConsentGranted())
		require.Equal(t, req.ExpiresAt.Unix(), got.ExpiresAt.Unix())
		got.ExpiresAt = req.ExpiresAt
		assert.EqualValues(t, req, got)

		require.NoError(t, m.AcceptConsentRequest(req.ID, new(oauth2_2.AcceptConsentRequestPayload)))
		got, err = m.GetConsentRequest(req.ID)
		require.NoError(t, err)
		assert.True(t, got.IsConsentGranted())

		require.NoError(t, m.RejectConsentRequest(req.ID, new(oauth2_2.RejectConsentRequestPayload)))
		got, err = m.GetConsentRequest(req.ID)
		require.NoError(t, err)
		assert.False(t, got.IsConsentGranted())
	})
}
