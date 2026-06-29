---
title: Time-slicing
linkTitle: Time-slicing
weight: 40
aliases:
  - /docs/guides/time-slicing/
description: Share a single GPU between multiple containers using CUDA time-slicing.
---

Time-slicing lets multiple containers share one physical GPU by scheduling their
CUDA contexts in turns. CUDA preempts each context at a configurable interval
and switches to the next one, similar to CPU time-sharing on a single-core
machine.

Use time-slicing when you want to increase GPU utilization across workloads
that do not run continuously — for example, multiple batch jobs or development
containers that would otherwise sit idle waiting for their own dedicated GPU.
For how time-slicing compares to MIG and the other sharing modes, refer to
[Choosing a resource type](../../concepts/gpu-allocation.md#choosing-a-resource-type).

## Feature status

`TimeSlicingSettings` is an Alpha feature gate, disabled by default.

| Feature gate | Default | Stage |
|---|---|---|
| `TimeSlicingSettings` | `false` | Alpha |

Refer to the [feature gates reference](../../reference/feature-gates.md) for all available gates.

## Prerequisites

- The DRA Driver for NVIDIA GPUs must be installed. Refer to [Installation](../../install.md).
- The `TimeSlicingSettings` feature gate must be enabled. Refer to [Enabling the feature](#enabling-the-feature).

## Enabling the feature

Enable the `TimeSlicingSettings` feature gate with `helm upgrade`:

```bash
helm upgrade dra-driver-nvidia-gpu oci://registry.k8s.io/dra-driver-nvidia/charts/dra-driver-nvidia-gpu \
  --namespace dra-driver-nvidia-gpu \
  --set featureGates.TimeSlicingSettings=true \
  --set gpuResourcesEnabledOverride=true
```

The GPU kubelet plugin and webhook must both restart for the change to take effect.
The rolling update happens automatically when you upgrade the Helm release.

## Configure time-slicing

Configuring time-slicing requires two steps:

1. Create a `ResourceClaimTemplate` that specifies the time-slice interval.
2. Create pods that reference the template.

A `ResourceClaimTemplate` is namespace-scoped. Create one in each namespace
where you want to use time-slicing. Within a namespace, create a separate
template for each interval configuration you need — for example, one template
for `interval: Long` batch workloads and another for `interval: Short`
interactive workloads. Pods that need the same interval share the same template.

### Time-slice intervals

The `interval` field controls the CUDA time-slice duration:

| Value | Description |
|---|---|
| `Default` | Uses the NVIDIA GPU driver's built-in default interval |
| `Short` | Shorter interval; contexts switch more frequently |
| `Medium` | Intermediate interval |
| `Long` | Longer interval; each context runs longer per turn before being preempted |

If omitted, `Default` is used.

### CEL selectors

You can add a `selectors` block under `exactly` to target a specific GPU by
model, UUID, or other attribute. Refer to
[Request full GPUs](allocating-gpus.md#select-a-gpu-by-product-name) for CEL selector examples
and available device attributes.

## Time-slicing example

### Create a ResourceClaimTemplate

A `ResourceClaimTemplate` defines the GPU request and its configuration.
Multiple pods can reuse the same template: Kubernetes automatically creates one
`ResourceClaim` per pod from it, and deletes that claim when the pod terminates.

1. Create a file called `shared-gpu.yaml`:

   ```yaml
   apiVersion: resource.k8s.io/v1
   kind: ResourceClaimTemplate
   metadata:
     namespace: time-slicing-example
     name: shared-gpu
   spec:
     spec:
       devices:
         requests:
         - name: gpu
           exactly:
             deviceClassName: gpu.nvidia.com
         config:
         - requests: ["gpu"]
           opaque:
             driver: gpu.nvidia.com
             parameters:
               apiVersion: resource.nvidia.com/v1beta1
               kind: GpuConfig
               sharing:
                 strategy: TimeSlicing
                 timeSlicingConfig:
                   interval: Long
   ```

   The `deviceClassName: gpu.nvidia.com` is required — it selects a full GPU.
   Set `interval` to the time-slice duration you need. Refer to
   [Time-slice intervals](#time-slice-intervals) for the available values.

2. Apply the manifest:

   ```bash
   kubectl apply -f shared-gpu.yaml
   ```

   Example output:

   ```
   resourceclaimtemplate.resource.k8s.io/shared-gpu created
   ```

### Create a Pod that references the ResourceClaimTemplate

Reference the `ResourceClaimTemplate` by name in `pod.spec.resourceClaims`.
Kubernetes creates one `ResourceClaim` per pod when it is scheduled.

To share a single GPU across containers in the same pod, each container
references the **same request name** (`request: gpu`) within the claim. If
containers referenced different request names, each would receive a separate
GPU.

1. Create a file called `time-slicing-pod.yaml`:

   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     namespace: time-slicing-example
     name: time-slicing-pod
   spec:
     containers:
     - name: workload-0
       image: <your-image>
       resources:
         claims:
         - name: shared-gpu
           request: gpu
     - name: workload-1
       image: <your-image>
       resources:
         claims:
         - name: shared-gpu
           request: gpu
     resourceClaims:
     - name: shared-gpu
       resourceClaimTemplateName: shared-gpu
     tolerations:
     - key: "nvidia.com/gpu"
       operator: "Exists"
       effect: "NoSchedule"
   ```

   Key fields:

   - **`image`** — replace `<your-image>` with your workload container image.
   - **`resourceClaimTemplateName`** — must match the name of the
     `ResourceClaimTemplate` you created in the previous step.
   - **`request: gpu`** — must match the request name defined in the template.
     Both containers using the same value is what causes them to share one GPU.
   - **Toleration** — allows the pod to schedule on nodes that have the
     `nvidia.com/gpu: NoSchedule` taint, which is common on GPU nodes. Remove
     it if your cluster does not use this taint.

2. Apply the manifest:

   ```bash
   kubectl apply -f time-slicing-pod.yaml
   ```

   Example output:

   ```
   pod/time-slicing-pod created
   ```

3. Verify the pod is running and both containers are ready:

   ```bash
   kubectl get pod -n time-slicing-example time-slicing-pod
   ```

   Example output:

   ```
   NAME                READY   STATUS    RESTARTS   AGE
   time-slicing-pod    2/2     Running   0          30s
   ```

4. Confirm both containers see the same GPU:

   ```bash
   kubectl exec -n time-slicing-example time-slicing-pod -c workload-0 -- nvidia-smi -L
   kubectl exec -n time-slicing-example time-slicing-pod -c workload-1 -- nvidia-smi -L
   ```

   Example output:

   ```
   GPU 0: NVIDIA A100-PCIE-40GB (UUID: GPU-2fa81118-5a5f-aa66-7660-471eed407181)
   GPU 0: NVIDIA A100-PCIE-40GB (UUID: GPU-2fa81118-5a5f-aa66-7660-471eed407181)
   ```

   Both commands return the same GPU UUID, confirming the containers share one device.

For additional examples, including time-slicing with CEL selectors, refer to the
[`demo/specs/`](https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu/tree/{{< param driver_release_tag >}}/demo/specs/)
directory in the repository.

## Limitations and considerations

- **No memory isolation.** All containers sharing a GPU access the same GPU
  memory. A container that allocates more memory than expected can affect
  other containers on the same device.
- **No throughput guarantees.** The GPU is shared on a best-effort basis.
  Workloads can observe variable performance depending on what else is running
  on the same GPU.
- **Not supported on MIG slices.** Setting `strategy: TimeSlicing` in a
  `MigDeviceConfig` is accepted without error but has no effect on hardware.
  To share a MIG slice across containers, use MPS instead.
- **Mutually exclusive with MPS on the same GPU.** Time-slicing sets the GPU
  compute mode to DEFAULT; MPS requires EXCLUSIVE_PROCESS. Both strategies can
  be active in a cluster, but not on the same physical GPU at the same time.
