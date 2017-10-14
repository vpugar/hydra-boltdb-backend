package internal

import (
	"github.com/gogo/protobuf/proto"
	"github.com/ory/hydra/warden/group"
)

//go:generate protoc --gofast_out=. group_record.proto

func GroupMarshal(v *group.Group) ([]byte, error) {
	return proto.Marshal(&GroupRecord{
		ID:      v.ID,
		Members: v.Members,
	})
}

func GroupUnmarshal(data []byte, v *group.Group) error {
	var pb GroupRecord
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	v.ID = pb.ID
	v.Members = pb.Members

	return nil
}
