package main

import (
	"encoding/json"
	"fmt"

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
	// If we can unmarshal to a Checkpoint directly we are done
	if err := json.Unmarshal(data, cp); err == nil {
		return nil
	}

	// Otherwise attempt to unmarshal to a Checkpoint2503RC2
	// TODO: Remove this one release cycle following the v25.3.0 release
	var cp2503rc2 Checkpoint2503RC2
	if err := cp2503rc2.UnmarshalCheckpoint(data); err != nil {
		return fmt.Errorf("unable to unmarshal as %T or %T", *cp, cp2503rc2)
	}

	// If that succeeded, verify its checksum...
	if err := cp2503rc2.VerifyChecksum(); err != nil {
		return fmt.Errorf("error verifying checksum for %T: %w", cp2503rc2, err)
	}

	// And then convert it to the v1 checkpoint format and try again
	cpconv := cp2503rc2.ToV1()
	data, err := cpconv.MarshalCheckpoint()
	if err != nil {
		return fmt.Errorf("error remarshalling %T after conversion from %T: %w", cpconv, cp2503rc2, err)
	}

	return cp.UnmarshalCheckpoint(data)
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
