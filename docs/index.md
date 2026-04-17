# Documentation

The DRA Driver for NVIDIA GPUs is a Kubernetes Dynamic Resource Allocation driver that lets workloads request GPUs, MIG slices, and VFIO passthrough devices (`gpu.nvidia.com`), and provisions ephemeral Multi-Node NVLink fabrics via IMEX (`compute-domain.nvidia.com`) so pods on different nodes can share GPU memory over NVLink. It targets Kubernetes 1.32+ with DRA enabled.

- [`architecture.md`](architecture.md) — components and request flows.
- [`development.md`](development.md) — build, test, and local-cluster loop.
