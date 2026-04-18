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

# install-gpu-operator.sh -- GPU Operator in minimal mode on the VM.
#
# Disables driver/toolkit/devicePlugin (the host + nvkind handle those),
# enables CDI, keeps NFD on for node labels. Waits for ClusterPolicy
# state=ready, since helm --wait only tracks pod readiness.

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# setup-nvkind-node.sh leaves tools under /usr/local/go and $HOME/go/bin,
# which a fresh SSH session does not pick up.
export PATH="/usr/local/go/bin:${HOME}/go/bin:${PATH}"

: "${NVKIND_CLUSTER_NAME:=dra-ci}"
: "${GPU_OPERATOR_VERSION:=v26.3.1}"

CTX="kind-${NVKIND_CLUSTER_NAME}"

helm repo add nvidia https://helm.ngc.nvidia.com/nvidia 2>/dev/null || true
helm repo update

helm --kube-context "${CTX}" upgrade -i gpu-operator nvidia/gpu-operator \
  --version "${GPU_OPERATOR_VERSION}" \
  --create-namespace -n gpu-operator \
  --set driver.enabled=false \
  --set toolkit.enabled=false \
  --set devicePlugin.enabled=false \
  --set cdi.enabled=true \
  --set cdi.default=true \
  --set mig.strategy=none \
  --set nfd.enabled=true

for i in $(seq 1 60); do
  state=$(kubectl --context "${CTX}" get clusterpolicy -o jsonpath='{.items[0].status.state}' 2>/dev/null || echo "")
  echo "[$((i*5))s] ClusterPolicy state: ${state:-<not-yet-created>}"
  [ "${state}" = "ready" ] && exit 0
  sleep 5
done

echo "ERROR: ClusterPolicy did not reach ready within 5 min" >&2
kubectl --context "${CTX}" get pods -n gpu-operator -o wide
kubectl --context "${CTX}" get clusterpolicy -o yaml | tail -40
exit 1
