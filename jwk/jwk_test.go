package jwk_test

import (
	"flag"
	"fmt"
	"github.com/ory/hydra/jwk"
	"github.com/vpugar/boltdbclient"
	jwk2 "github.com/vpugar/hydra-boltdb-backend/jwk"
	"os"
	"testing"
)

const (
	name   = "bolt"
	dbName = "jwk2.test.db"
)

var jwkManager jwk.Manager

var testGenerator = &jwk.RS256Generator{}

func startTest(m *testing.M) int {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	if cm, err := jwk2.NewJwkManager(boltClient); err != nil {
		os.Exit(1)
	} else {
		jwkManager = cm
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
	var _ jwk.Manager = &jwk2.JwkManager{}
}

func TestManagerKey(t *testing.T) {
	ks, _ := testGenerator.Generate("")

	t.Run(name, func(t *testing.T) {
		jwk.TestHelperManagerKey(jwkManager, ks)(t)
	})
}

func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	ks.Key("private")

	t.Run(name, func(t *testing.T) {
		jwk.TestHelperManagerKeySet(jwkManager, ks)(t)
	})
}
