#!/usr/bin/env bash

# Copyright The Kubernetes Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o pipefail
set -o errexit
set -o nounset

# Debug-log input values for visibility in CI logs, these are expected to be set by BuildKit.
declare -p TARGETARCH
declare -p BASH_STATIC_GIT_REF

# Dependencies for bash-static build.
apt-get update
apt-get install -y gpg curl autoconf file git

git clone https://github.com/robxu9/bash-static/

# Map BuildKit TARGETARCH to bash-static arch names, always set CC/STRIP (required with set -u).
ARCH="${TARGETARCH}"
case "${ARCH}" in
  arm64) ARCH=aarch64 ;;
  amd64) ARCH=x86_64 ;;
esac
STRIP=strip
CC=gcc
case "${ARCH}" in
  aarch64)
    STRIP=aarch64-linux-gnu-strip
    CC=aarch64-linux-gnu-gcc
    ;;
esac

echo "detected arch: ${ARCH}"
echo "cc to use: ${CC}"
echo "strip to use: ${STRIP}"

cd bash-static
git checkout "${BASH_STATIC_GIT_REF}"
sed -i 's|https://ftp\.gnu\.org/gnu|https://mirrors.kernel.org/gnu/|g' ./build.sh
sed -i 's/-sLO/-sSfLO --retry 300 --connect-timeout 20 --retry-delay 5 --retry-all-errors /g' ./build.sh
sed -i "s/strip/${STRIP}/g" ./build.sh
sed -i 's/make -s \&\& make -s tests/make -j4/g' ./build.sh
bash version-52.sh
STRIP="${STRIP}" CC="${CC}" ./build.sh linux "${ARCH}"
