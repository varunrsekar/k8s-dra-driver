/*
Copyright The Kubernetes Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package imex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		gateEnabled bool
		expectError bool
		description string
	}{
		{
			name:        "driverManaged with empty isolation defaults to IMEXDomain",
			config:      Config{Mode: ModeDriverManaged, Isolation: ""},
			gateEnabled: false,
			expectError: false,
			description: "empty isolation is treated the same as the IMEXDomain default",
		},
		{
			name:        "driverManaged with explicit IMEXDomain isolation succeeds",
			config:      Config{Mode: ModeDriverManaged, Isolation: IsolationIMEXDomain},
			gateEnabled: false,
			expectError: false,
			description: "IMEXDomain is the default, supported isolation strategy under driverManaged too",
		},
		{
			name:        "driverManaged with IMEXChannel isolation errors",
			config:      Config{Mode: ModeDriverManaged, Isolation: IsolationIMEXChannel},
			gateEnabled: false,
			expectError: true,
			description: "per-workload channel isolation is not implemented yet, regardless of mode",
		},
		{
			name:        "driverManaged with unknown isolation errors",
			config:      Config{Mode: ModeDriverManaged, Isolation: "bogus"},
			gateEnabled: true,
			expectError: true,
			description: "isolation is validated regardless of mode",
		},
		{
			name:        "hostManaged requires the feature gate",
			config:      Config{Mode: ModeHostManaged, Isolation: IsolationIMEXDomain},
			gateEnabled: false,
			expectError: true,
			description: "mode=hostManaged is rejected when HostManagedIMEXDaemon is disabled",
		},
		{
			name:        "hostManaged with empty isolation defaults to IMEXDomain",
			config:      Config{Mode: ModeHostManaged, Isolation: ""},
			gateEnabled: true,
			expectError: false,
			description: "empty isolation is treated the same as the IMEXDomain default under hostManaged too",
		},
		{
			name:        "hostManaged with explicit IMEXDomain isolation succeeds",
			config:      Config{Mode: ModeHostManaged, Isolation: IsolationIMEXDomain},
			gateEnabled: true,
			expectError: false,
			description: "IMEXDomain is the default isolation strategy, supported under hostManaged",
		},
		{
			name:        "hostManaged with IMEXChannel isolation errors for now",
			config:      Config{Mode: ModeHostManaged, Isolation: IsolationIMEXChannel},
			gateEnabled: true,
			expectError: true,
			description: "per-workload channel isolation is not implemented yet, so it errors out",
		},
		{
			name:        "unknown mode errors",
			config:      Config{Mode: "bogus"},
			gateEnabled: true,
			expectError: true,
			description: "an unrecognized mode value must be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(tt.gateEnabled)
			if tt.expectError {
				require.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
			}
		})
	}
}

func TestConfigEffectiveHostManaged(t *testing.T) {
	require.True(t, Config{Mode: ModeHostManaged}.EffectiveHostManaged())
	require.False(t, Config{Mode: ModeDriverManaged}.EffectiveHostManaged())
	require.False(t, Config{}.EffectiveHostManaged())
}
