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

DRIVER_NAME := nvidia-dra-driver-gpu
MODULE := sigs.k8s.io/$(DRIVER_NAME)

REGISTRY ?= registry.k8s.io/nv-dra-driver-gpu
# Staging registry before kubernetes/k8s.io promotion (Prow / test-infra jobs).
# See cloudbuild.yaml and https://github.com/kubernetes/k8s.io (k8s-staging-nv-dra-driver-gpu).
STAGING_REGISTRY ?= us-central1-docker.pkg.dev/k8s-staging-images/nv-dra-driver-gpu

VERSION  ?= v26.4.0-dev

# vVERSION represents the version with a guaranteed v-prefix
# Note: this is probably not consumed in our build chain.
# `VERSION` above is expected to have a `v` prefix, which is
# then automatically stripped in places that must not have it
# (e.g., in context of Helm).
vVERSION := v$(VERSION:v%=%)

GOLANG_VERSION := $(shell ./hack/golang-version.sh)
TOOLKIT_CONTAINER_IMAGE := $(shell ./hack/toolkit-container-image.sh)
BASH_STATIC_GIT_REF := 021f5f29f665c92ca16a369d9f27e288c3aed0c6

# These variables are only needed when building a local image
BUILDIMAGE_TAG ?= devel-go$(GOLANG_VERSION)
BUILDIMAGE ?=  $(DRIVER_NAME):$(BUILDIMAGE_TAG)

GIT_COMMIT ?= $(shell git describe --match="" --dirty --long --always --abbrev=40 2> /dev/null || echo "")
GIT_COMMIT_SHORT ?= $(shell git rev-parse --short=8 HEAD)

# Shape: v25.8.0-dev-f2eaddd6
VERSION_W_COMMIT = $(VERSION)-$(GIT_COMMIT_SHORT)

# Shape: 25.8.0-dev-f2eaddd6-chart (no leading v)
VERSION_STAGING_CHART ?= $(VERSION_W_COMMIT:v%=%)-chart
VERSION_GHCR_CHART ?= $(VERSION_STAGING_CHART)

print-%:
	@echo $($*)
