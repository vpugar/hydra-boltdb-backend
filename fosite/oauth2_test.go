package fosite_test

import (
	"flag"
	"fmt"
	"github.com/ory/fosite"
	client2 "github.com/ory/hydra/client"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/vpugar/boltdbclient"
	"github.com/vpugar/hydra-boltdb-backend/client"
	fosite2 "github.com/vpugar/hydra-boltdb-backend/fosite"
	"os"
	"testing"
)

const (
	name   = "bolt"
	dbName = "oauth2.test.db"
)

var clientManager pkg.FositeStorer

func startTest(m *testing.M) int {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	var hasher fosite.Hasher = &fosite.BCrypt{}
	if cm, err := client.NewClientManager(boltClient, &hasher); err != nil {
		os.Exit(1)
	} else {

		cm.CreateClient(&client2.Client{
			ID: "foobar",
		})

		if om, err := fosite2.NewOauth2Manager(boltClient, cm); err != nil {
			os.Exit(1)
		} else {
			clientManager = om
		}
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
	var _ pkg.FositeStorer = &fosite2.Oauth2Manager{}
}

func TestCreateGetDeleteAuthorizeCodes(t *testing.T) {
	t.Run(name, oauth2.TestHelperCreateGetDeleteAuthorizeCodes(clientManager))

}

func TestCreateGetDeleteAccessTokenSession(t *testing.T) {
	t.Run(name, oauth2.TestHelperCreateGetDeleteAccessTokenSession(clientManager))
}

func TestCreateGetDeleteOpenIDConnectSession(t *testing.T) {
	t.Run(name, oauth2.TestHelperCreateGetDeleteOpenIDConnectSession(clientManager))
}

func TestCreateGetDeleteRefreshTokenSession(t *testing.T) {
	t.Run(name, oauth2.TestHelperCreateGetDeleteRefreshTokenSession(clientManager))
}

func TestRevokeRefreshToken(t *testing.T) {
	t.Run(name, oauth2.TestHelperRevokeRefreshToken(clientManager))
}
