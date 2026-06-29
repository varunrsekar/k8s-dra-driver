---
title: View available GPU resources
linkTitle: View GPU resources
weight: 20
aliases:
  - /docs/guides/gpu-allocation/view-resourceslices/
description: Inspect ResourceSlices and pools to see which GPU, MIG, and VFIO devices are allocatable on your cluster.
---

This guide shows how to inspect the `ResourceSlice` objects and pools the DRA
Driver for NVIDIA GPUs publishes, so you can see which GPU, MIG, and VFIO devices
are allocatable on your cluster and confirm each node is advertising the devices
you expect.

For what each device attribute and capacity field means and the exact names to
use in CEL selectors, refer to
[ResourceSlice device attributes](../../reference/resourceslice-attributes.md).
For details on how the driver publishes ResourceSlices, refer to
[Publishing GPUs in ResourceSlices](../../concepts/gpu-allocation.md#publishing-gpus-in-resourceslices).
For the generic DRA model (what ResourceSlices, DeviceClasses, and ResourceClaims
are and how the scheduler uses them), refer to the
[Kubernetes DRA documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/).

## Prerequisites

- The DRA Driver for NVIDIA GPUs is installed. Refer to [Installation](../../install.md).
- At least one GPU node with devices published in a `ResourceSlice`.
- [`jq`](https://jqlang.github.io/jq/) installed, for the JSON filtering example below (optional).

## List ResourceSlices

List all ResourceSlices in the cluster:

```bash
kubectl get resourceslice
```

Example output:

```
NAME                                        NODE           DRIVER           POOL           AGE
worker-gpu-01-gpu.nvidia.com-abc12          worker-gpu-01  gpu.nvidia.com   worker-gpu-01  3m
```

If this list is empty, no node is advertising GPU devices yet. Confirm the driver
is installed and running and that at least one node has supported GPUs.

## Inspect a slice in detail

Inspect the full detail of a slice, including per-device attributes and capacity:

```bash
kubectl get resourceslice -o yaml
```

This returns a `List` whose `items` are `ResourceSlice` objects — one pool per
node. Each slice wraps the node's allocatable devices under `spec.devices`:

```yaml
apiVersion: v1
kind: List
items:
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    name: 00000-gpu.nvidia.com-ipp1-0033-4lpks
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: ipp1-0033            # the node whose devices this slice describes
      uid: 2ed85ae9-1a62-42e6-a6ac-1a6ef5f787c1
    # creationTimestamp, generateName, generation, resourceVersion, uid also present
  spec:
    driver: gpu.nvidia.com        # always gpu.nvidia.com for GPU, MIG, and VFIO
    nodeName: ipp1-0033
    pool:
      name: ipp1-0033
      generation: 1
      resourceSliceCount: 1
    devices:
    - ...                         # one entry per allocatable device
```

Each `spec.devices[]` entry describes one device. For a representative entry for
each device type (full GPU, MIG slice, VFIO passthrough) and what every attribute
and capacity field means, refer to
[ResourceSlice device attributes](../../reference/resourceslice-attributes.md).

## List only NVIDIA GPU devices

All GPU, MIG, and VFIO devices are published by the same driver,
`gpu.nvidia.com`. To list just those devices with their attributes and capacity:

```bash
kubectl get resourceslice -o json \
  | jq -r '.items[]
    | select(.spec.driver == "gpu.nvidia.com")
    | .spec.devices[]
    | {name: .name, type: .attributes.type.string,
       attributes: .attributes, capacity: .capacity}'
```

ComputeDomain IMEX channels are published separately under the
`compute-domain.nvidia.com` driver. Refer to
[ComputeDomain workloads](../compute-domain-workloads.md) for details.

## Check resource pools

The driver publishes one pool per node, named after the node. A pool groups all
of a node's ResourceSlices together, so the scheduler treats them as a single,
consistent view of that node's devices. To review the pools across your cluster:

```bash
kubectl get resourceslice -o custom-columns=\
POOL:.spec.pool.name,\
NODE:.spec.nodeName,\
SLICES:.spec.pool.resourceSliceCount,\
GENERATION:.spec.pool.generation
```

Example output:

```
POOL            NODE            SLICES   GENERATION
worker-gpu-01   worker-gpu-01   1        1
worker-gpu-02   worker-gpu-02   1        1
```

- `SLICES` (`spec.pool.resourceSliceCount`) is how many ResourceSlices make up
  the pool. This is `1` for most nodes; a node with many devices can be split
  across several slices. The scheduler only allocates from a pool once it sees
  all of its slices at the same generation.
- `GENERATION` (`spec.pool.generation`) increases each time the driver
  republishes the pool — for example after devices change or the GPU kubelet
  plugin restarts. All slices in a pool share the current generation.

If a node with GPUs has no pool, or its slice count is lower than the number of
GPU nodes you expect, the GPU kubelet plugin on that node may not be running.

## What's next

Start requesting GPUs in your workloads:

- [Request full GPUs](allocating-gpus.md)
- [Time-slicing](time-slicing.md)
- [MIG](mig.md)
- [VFIO passthrough](kubevirt-vfio-gpu-passthrough.md)