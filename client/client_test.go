package client_test

import (
	"flag"
	"fmt"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/vpugar/boltdbclient"
	client2 "github.com/vpugar/hydra-boltdb-backend/client"
	"os"
	"testing"
)

const (
	name   = "bolt"
	dbName = "client.test.db"
)

var clientManager client.Manager

func startTest(m *testing.M) int {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	var hasher fosite.Hasher = &fosite.BCrypt{}

	if cm, err := client2.NewClientManager(boltClient, hasher); err != nil {
		os.Exit(1)
	} else {
		clientManager = cm
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
	var _ client.Manager = &client2.ClientManager{}
}

func TestCreateGetDeleteClient(t *testing.T) {
	t.Run(name, client.TestHelperCreateGetDeleteClient(name, clientManager))
}

func TestClientAutoGenerateKey(t *testing.T) {
	t.Run(name, client.TestHelperClientAutoGenerateKey(name, clientManager))
}

func TestAuthenticateClient(t *testing.T) {
	t.Run(name, client.TestHelperClientAuthenticate(name, clientManager))
}
