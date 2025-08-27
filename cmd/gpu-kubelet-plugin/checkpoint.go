package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/checksum"
)

type Checkpoint struct {
	Checksum checksum.Checksum `json:"checksum"`
	V1       *CheckpointV1     `json:"v1,omitempty"`
	V2       *CheckpointV2     `json:"v2,omitempty"`
}

func (cp *Checkpoint) ToLatestVersion() *Checkpoint {
	latest := &Checkpoint{}
	switch {
	case cp.V2 != nil:
		latest.V2 = cp.V2
	case cp.V1 != nil:
		latest.V2 = cp.V1.ToV2()
	default:
		latest.V2 = &CheckpointV2{}
	}
	if latest.V2.PreparedClaims == nil {
		latest.V2.PreparedClaims = make(PreparedClaimsByUID)
	}
	return latest
}

func (cp *Checkpoint) MarshalCheckpoint() ([]byte, error) {
	cp.Checksum = 0
	out, err := json.Marshal(*cp)
	if err != nil {
		return nil, err
	}
	cp.Checksum = checksum.New(out)
	return json.Marshal(*cp)
}

func (cp *Checkpoint) UnmarshalCheckpoint(data []byte) error {
	return json.Unmarshal(data, cp)
}

func (cp *Checkpoint) VerifyChecksum() error {
	ck := cp.Checksum
	cp.Checksum = 0
	defer func() {
		cp.Checksum = ck
	}()
	out, err := json.Marshal(*cp)
	if err != nil {
		return err
	}
	return ck.Verify(out)
}
