/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"testing"

	"github.com/NVIDIA/go-nvlib/pkg/nvpci"
	"github.com/stretchr/testify/require"
)

func TestDiscoverNumaNodeUsesPCI(t *testing.T) {
	pci := &nvpci.InterfaceMock{
		GetGPUByPciBusIDFunc: func(pciBusID string) (*nvpci.NvidiaPCIDevice, error) {
			require.Equal(t, "0000:9f:00.0", pciBusID)
			return &nvpci.NvidiaPCIDevice{NumaNode: 2}, nil
		},
	}
	lib := deviceLib{nvpci: pci}

	numaNode, err := lib.discoverNumaNode("0000:9f:00.0")

	require.NoError(t, err)
	require.NotNil(t, numaNode)
	require.Equal(t, 2, *numaNode)
	require.Len(t, pci.GetGPUByPciBusIDCalls(), 1)
}

func TestDiscoverNumaNodeOmitsUnknownPCINumaNode(t *testing.T) {
	pci := &nvpci.InterfaceMock{
		GetGPUByPciBusIDFunc: func(string) (*nvpci.NvidiaPCIDevice, error) {
			return &nvpci.NvidiaPCIDevice{NumaNode: -1}, nil
		},
	}
	lib := deviceLib{nvpci: pci}

	numaNode, err := lib.discoverNumaNode("0000:9f:00.0")

	require.NoError(t, err)
	require.Nil(t, numaNode)
}
