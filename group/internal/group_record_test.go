package internal_test

import (
	"github.com/ory/hydra/warden/group"
	"github.com/vpugar/hydra-boltdb-backend/group/internal"
	"reflect"
	"testing"
)

func TestMarshalGroup(t *testing.T) {
	v := group.Group{
		ID:      "ID",
		Members: []string{"1", "2"},
	}

	var other group.Group
	if buf, err := internal.GroupMarshal(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.GroupUnmarshal(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
