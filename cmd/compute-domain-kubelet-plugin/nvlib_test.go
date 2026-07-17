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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCheckHostIMEXReadyMissingBinary confirms that a driver root without
// nvidia-imex-ctl (e.g. a node where the optional nvidia-imex package was
// never installed) fails with an actionable error, without ever attempting
// to chroot/exec. The success path (chrooting in and parsing "READY") needs
// a real nvidia-imex-ctl binary and elevated privileges, so it is not
// covered by this unit test.
func TestCheckHostIMEXReadyMissingBinary(t *testing.T) {
	l := deviceLib{devRoot: t.TempDir()}

	err := l.checkHostIMEXReady(defaultIMEXHostSocketPath)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "nvidia-imex-ctl not found")
}
