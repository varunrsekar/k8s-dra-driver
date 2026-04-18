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

# collect-artifacts.sh -- run on the VM to dump cluster state for postmortem.
# If ARTIFACTS is set, also writes kind's per-node logs alongside.

export PATH="/usr/local/go/bin:${HOME}/go/bin:${PATH}"
CTX="kind-${NVKIND_CLUSTER_NAME:-dra-ci}"

echo "=== nvidia-smi ==="
nvidia-smi 2>&1 || true
echo
echo "=== driver version ==="
cat /sys/module/nvidia/version 2>&1 || true
echo
echo "=== nodes ==="
kubectl --context "${CTX}" get nodes -o wide 2>&1 || true
echo
echo "=== pods (all namespaces) ==="
kubectl --context "${CTX}" get pods -A -o wide 2>&1 || true
echo
echo "=== pod describe (driver namespace) ==="
kubectl --context "${CTX}" describe pods -n dra-driver-nvidia-gpu 2>&1 || true
echo
echo "=== pod describe (gpu-operator namespace) ==="
kubectl --context "${CTX}" describe pods -n gpu-operator 2>&1 || true
echo
echo "=== resourceslices ==="
kubectl --context "${CTX}" get resourceslices -A -o yaml 2>&1 || true
echo
echo "=== deviceclasses ==="
kubectl --context "${CTX}" get deviceclass -o yaml 2>&1 || true
echo
echo "=== resourceclaims ==="
kubectl --context "${CTX}" get resourceclaims -A -o yaml 2>&1 || true
echo
echo "=== events ==="
kubectl --context "${CTX}" get events -A --sort-by=.lastTimestamp 2>&1 || true
echo
echo "=== dra driver kubelet-plugin logs (gpus) ==="
kubectl --context "${CTX}" logs -n dra-driver-nvidia-gpu -l dra-driver-nvidia-gpu-component=kubelet-plugin -c gpus --tail=500 2>&1 || true
echo
echo "=== dra driver controller logs ==="
kubectl --context "${CTX}" logs -n dra-driver-nvidia-gpu -l dra-driver-nvidia-gpu-component=controller --tail=200 2>&1 || true
echo
echo "=== clusterpolicy ==="
kubectl --context "${CTX}" get clusterpolicy -o yaml 2>&1 || true

# kind writes logs (kubelet, containerd, etcd, apiserver) per-node into a
# directory. If ARTIFACTS is set, drop them there alongside the stdout dump.
if [ -n "${ARTIFACTS:-}" ]; then
  echo "=== exporting kind cluster logs to ${ARTIFACTS}/kind-logs ==="
  sg docker -c "kind export logs '${ARTIFACTS}/kind-logs' --name '${NVKIND_CLUSTER_NAME:-dra-ci}'" 2>&1 || true
fi
