package group_test

import (
	"flag"
	"fmt"
	"github.com/ory/hydra/warden/group"
	"github.com/vpugar/boltdbclient"
	group2 "github.com/vpugar/hydra-boltdb-backend/group"
	"os"
	"testing"
)

const (
	name   = "bolt"
	dbName = "group.test.db"
)

var groupManager group.Manager

func startTest(m *testing.M) int {

	boltClient := boltdbclient.NewClient(boltdbclient.Config{
		Dir:      "./",
		Filename: dbName,
	})
	boltClient.Open()
	defer stopTest(boltClient)

	if cm, err := group2.NewGroupManager(boltClient); err != nil {
		os.Exit(1)
	} else {
		groupManager = cm
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
	var _ group.Manager = &group2.GroupManager{}
}

func TestManagers(t *testing.T) {
	t.Run(name, group.TestHelperManagers(groupManager))
}
