package ladon_test

import (
	"flag"
	"fmt"
	"github.com/ory/ladon"
	"github.com/vpugar/boltdbclient"
	ladon2 "github.com/vpugar/hydra-boltdb-backend/ladon"
	"os"
	"testing"
)

const (
	name   = "bolt"
	dbName = "ladon.test.db"
)

var ladonManager ladon.Manager

func startTest(m *testing.M) int {
	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	if cm, err := ladon2.NewLadonManager(boltClient); err != nil {
		fmt.Println("FAIL")
		return 1
	} else {
		ladonManager = cm
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
	var _ ladon.Manager = &ladon2.LadonManager{}
}

func TestGetErrors(t *testing.T) {
	t.Run(name, ladon.TestHelperGetErrors(ladonManager))

}

func TestCreateGetDelete(t *testing.T) {
	t.Run(name, ladon.TestHelperCreateGetDelete(ladonManager))
}

//func TestFindPoliciesForSubject(t *testing.T) {
// FIXME t.Run(name, ladon.TestHelperFindPoliciesForSubject(name, ladonManager))
//}
