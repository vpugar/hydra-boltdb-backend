package backend_test

import (
	"github.com/stretchr/testify/require"
	"github.com/vpugar/hydra-boltdb-backend/backend"
	"testing"
	//"os"
	"os"
)

func TestNewBoltdbConnection(t *testing.T) {
	c := backend.NewBoltdbConnection("backend_test1.db")
	require.NotEmpty(t, c, "No connection")

	err := c.Connect()
	require.NoError(t, err, "Cannot connect")

	err = c.Disconnect()
	require.NoError(t, err, "Cannot disconnect")

	os.Remove("./backend_test1.db")
	require.NoError(t, err, "Cannot delete db file")
}

func diconnect(t *testing.T, c *backend.BoltdbConnection) {
	err := c.Disconnect()
	require.NoError(t, err, "Cannot disconnect")

	os.Remove("./backend_test2.db")
	require.NoError(t, err, "Cannot delete db file")
}

func TestInterfacesConnection(t *testing.T) {
	c := backend.NewBoltdbConnection("backend_test2.db")
	require.NotEmpty(t, c, "No connection")

	err := c.Connect()
	require.NoError(t, err, "Cannot connect")

	defer diconnect(t, c)

	gm, err := c.NewGroupManager()
	require.NoError(t, err, "Cannot create group management")
	require.NotEmpty(t, gm, "No group management")

	pm, err := c.NewPolicyManager()
	require.NoError(t, err, "Cannot create policy management")
	require.NotEmpty(t, pm, "No policy management")

	cm, err := c.NewClientManager(nil)
	require.NoError(t, err, "Cannot create client management")
	require.NotEmpty(t, cm, "No client management")

	om, err := c.NewOAuth2Manager(cm)
	require.NoError(t, err, "Cannot create oauth2 management")
	require.NotEmpty(t, om, "No oauth2 management")

	jm, err := c.NewJWKManager(nil)
	require.NoError(t, err, "Cannot create jwk management")
	require.NotEmpty(t, jm, "No jwk management")

	crm, err := c.NewConsentRequestManager()
	require.NoError(t, err, "Cannot create consent request management")
	require.NotEmpty(t, crm, "No consent request management")

}
