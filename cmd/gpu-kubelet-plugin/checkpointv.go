package main

import (
	resourceapi "k8s.io/api/resource/v1"
)

type ClaimCheckpointState string

const (
	ClaimCheckpointStateUnset            ClaimCheckpointState = ""
	ClaimCheckpointStatePrepareCompleted ClaimCheckpointState = "PrepareCompleted"
)

// Latest version type aliases

type PreparedClaimsByUID = PreparedClaimsByUIDV2
type PreparedClaim = PreparedClaimV2

// V2 types

type CheckpointV2 struct {
	PreparedClaims PreparedClaimsByUIDV2 `json:"preparedClaims,omitempty"`
}

type PreparedClaimsByUIDV2 map[string]PreparedClaimV2

type PreparedClaimV2 struct {
	CheckpointState ClaimCheckpointState            `json:"checkpointState"`
	Status          resourceapi.ResourceClaimStatus `json:"status,omitempty"`
	PreparedDevices PreparedDevices                 `json:"preparedDevices,omitempty"`
}

// V1 types

type CheckpointV1 struct {
	PreparedClaims PreparedClaimsByUIDV1 `json:"preparedClaims,omitempty"`
}

type PreparedClaimsByUIDV1 map[string]PreparedClaimV1

type PreparedClaimV1 struct {
	Status          resourceapi.ResourceClaimStatus `json:"status,omitempty"`
	PreparedDevices PreparedDevices                 `json:"preparedDevices,omitempty"`
}

// Conversion functions

func (v1 *CheckpointV1) ToV2() *CheckpointV2 {
	v2 := &CheckpointV2{
		PreparedClaims: make(PreparedClaimsByUIDV2),
	}
	for claimUID, v1Claim := range v1.PreparedClaims {
		v2.PreparedClaims[claimUID] = PreparedClaimV2{
			CheckpointState: ClaimCheckpointStatePrepareCompleted,
			Status:          v1Claim.Status,
			PreparedDevices: v1Claim.PreparedDevices,
		}
	}
	return v2
}
