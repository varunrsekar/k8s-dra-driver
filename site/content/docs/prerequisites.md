---
title: Prerequisites
linkTitle: Prerequisites
weight: 30
description: Cluster, node, and tooling requirements before installing the driver.
---

Cluster, software, and hardware requirements for the DRA Driver for NVIDIA GPUs.

{{% alert title="Tip" %}}
Most of these prerequisites can be installed and managed for you by the [NVIDIA GPU Operator](#install-prerequisites-with-nvidia-gpu-operator).
{{% /alert %}}


| Requirement | Version / Notes |
|---|---|
| Kubernetes | v1.34.2 or later, with at least one node that has one or more NVIDIA GPUs. The use of DRA became GA in Kubernetes v1.34+ and earlier versions required the `DynamicResourceAllocation` feature gate. |
| Helm | v3.8 or later. |
| NVIDIA Driver | v565 or later for GPU allocation. v570.158.01 or later if using [ComputeDomains](#computedomains-additional-prerequisites). |
| CDI  | Enabled in your container runtime. This is enabled by default in containerd 2.0+ and CRIO v1.27+. The DRA Driver uses CDI to expose GPUs to containers.  |
| Node Feature Discovery (NFD) | Labels GPU nodes in the cluster. The DRA Driver uses these labels to target the GPU kubelet plugin to the correct nodes. |

## ComputeDomains additional prerequisites

If you plan to use ComputeDomains, you also need:

- NVIDIA Driver v570.158.01 or later. The `IMEXDaemonsWithDNSNames` feature gate is enabled by default and requires this driver version. The ComputeDomain plugin will fail to start on older drivers unless `IMEXDaemonsWithDNSNames` is explicitly disabled.
- Multi-Node NVLink (MNNVL) hardware. Nodes must be connected via NVLink fabric, such as GB200 NVL72 and similar systems.
- GPU Feature Discovery (GFD) deployed via the [GPU Operator](#install-prerequisites-with-nvidia-gpu-operator). GFD generates the `nvidia.com/gpu.clique` node labels required by ComputeDomains.
- On all GPU nodes where the `nvidia-imex-*` packages are installed, the `nvidia-imex.service` systemd unit must be disabled:

```bash
systemctl disable --now nvidia-imex.service && systemctl mask nvidia-imex.service
```

### Host-managed IMEX

By default the driver owns the `nvidia-imex` daemon lifecycle, per the requirement above. For clusters where the operator already runs `nvidia-imex` as a host service (for example through systemd), set `resources.computeDomains.imex.mode=hostManaged` to stop the driver from creating per-ComputeDomain `nvidia-imex` DaemonSets. This **inverts** the requirement above: `nvidia-imex.service` must be installed, configured, and left **running** on every participating GPU node.

Host-managed mode requires two Helm values together:

```bash
--set featureGates.HostManagedIMEXDaemon=true \
--set resources.computeDomains.imex.mode=hostManaged
```

- `featureGates.HostManagedIMEXDaemon` is an alpha gate that only unlocks setting `imex.mode=hostManaged`. It does not by itself change any behavior. Setting `imex.mode=hostManaged` without this gate enabled is a startup validation error (`helm install`/`helm template` fails immediately, and if bypassed, the controller and kubelet plugin pods will fail to start).
- `imex.isolation` selects the IMEX isolation strategy, and applies under both `driverManaged` and `hostManaged` mode. It must be set to one of:
  - `IMEXDomain` (default) — all workloads running in the same IMEX domain share the same channel (0). Under `driverManaged` mode this is inherent to the model (the driver creates one `nvidia-imex` daemon per ComputeDomain); under `hostManaged` mode, multiple ComputeDomains can still run against the same host IMEX domain, all receiving channel 0, with no isolation between them.
  - `IMEXChannel` — intended to eventually give each workload a unique channel within an IMEX domain. Not implemented yet: setting it is a startup validation error regardless of `imex.mode`.
  - Any other value is a startup validation error.
- Changing `imex.mode` or `imex.isolation` on a cluster with active ComputeDomain workloads is not supported and is not enforced by the driver, drain and remove existing ComputeDomains first.

## Install prerequisites with NVIDIA GPU Operator

The [NVIDIA GPU Operator](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/index.html) is a Kubernetes operator that automates the deployment and lifecycle management of all NVIDIA software components needed to provision and monitor GPUs in a cluster.

It can manage the following DRA Driver for NVIDIA GPUs prerequisites for you:

- NVIDIA Driver (v565+ for GPU allocation, v570.158.01+ for ComputeDomains). The GPU Operator installs a [default driver](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/platform-support.html#gpu-operator-component-matrix) that meets the DRA Driver's prerequisites. To use a specific version, see [Common chart customization options](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/getting-started.html#common-chart-customization-options) in the GPU Operator documentation.
- CDI enabled through the NVIDIA Container Toolkit.
- Node Feature Discovery (NFD).
- GPU Feature Discovery (GFD), required for ComputeDomains.

If you choose to install the GPU Operator, follow the [DRA Driver for NVIDIA GPUs install guide](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/dra-intro-install.html) in the GPU Operator documentation. It covers installing the GPU Operator with the NVIDIA Kubernetes Device Plugin disabled and installing the DRA Driver for NVIDIA GPUs.
