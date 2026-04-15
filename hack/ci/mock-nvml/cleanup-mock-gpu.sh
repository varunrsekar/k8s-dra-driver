#!/usr/bin/env bash
# Copyright The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# cleanup-mock-gpu.sh -- Remove mock NVML installation from the host.
#
# Reverses everything setup-mock-gpu.sh created. Safe to run multiple times.
set -o errexit
set -o nounset
set -o pipefail

MOCK_ROOT="${MOCK_ROOT:-/var/lib/nvml-mock}"

echo "=== Cleaning up mock GPU environment ==="

# Remove mock driver root (includes device nodes, config, libs)
if [ -d "${MOCK_ROOT}" ]; then
  echo "Removing ${MOCK_ROOT}/"
  sudo rm -rf "${MOCK_ROOT}"
fi

# Remove CDI spec
sudo rm -f /var/run/cdi/nvidia.yaml

# Remove GPU Operator compat symlink
sudo rm -f /run/nvidia/driver

echo "Cleanup complete"
