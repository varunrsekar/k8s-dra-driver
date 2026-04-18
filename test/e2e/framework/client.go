// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

// Package framework contains shared helpers for the dra-driver-nvidia-gpu
// e2e suite: client construction, ResourceSlice-based GPU detection,
// manifest rendering, and polling helpers.
package framework

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// BuildConfig returns an in-cluster *rest.Config if available, otherwise
// uses the standard KUBECONFIG / ~/.kube/config discovery.
func BuildConfig() (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}
	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig: %w", err)
	}
	return cfg, nil
}

// NewClientset builds a typed Kubernetes clientset; DRA v1 types live
// under Clientset().ResourceV1().
func NewClientset() (*kubernetes.Clientset, *rest.Config, error) {
	cfg, err := BuildConfig()
	if err != nil {
		return nil, nil, err
	}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("build clientset: %w", err)
	}
	return cs, cfg, nil
}
