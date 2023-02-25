/*
 * Copyright (c) 2022, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nascrd "github.com/NVIDIA/k8s-dra-driver/api/nvidia.com/resource/gpu/nas/v1alpha1"
	nvclientset "github.com/NVIDIA/k8s-dra-driver/pkg/nvidia.com/resource/clientset/versioned"
)

const (
	NodeAllocationStateStatusReady    = "Ready"
	NodeAllocationStateStatusNotReady = "NotReady"
)

type NodeAllocationStateConfig struct {
	Name      string
	Namespace string
	Owner     *metav1.OwnerReference
}

type MigDevicePlacement = nascrd.MigDevicePlacement
type AllocatableGpu = nascrd.AllocatableGpu
type AllocatableMigDevices = nascrd.AllocatableMigDevices
type AllocatableDevices = nascrd.AllocatableDevices
type AllocatedGpu = nascrd.AllocatedGpu
type AllocatedMigDevice = nascrd.AllocatedMigDevice
type AllocatedDevice = nascrd.AllocatedDevice
type AllocatedDevices = nascrd.AllocatedDevices
type RequestedGpu = nascrd.RequestedGpu
type RequestedMigDevice = nascrd.RequestedMigDevice
type RequestedGpus = nascrd.RequestedGpus
type RequestedMigDevices = nascrd.RequestedMigDevices
type RequestedDevices = nascrd.RequestedDevices
type NodeAllocationStateSpec = nascrd.NodeAllocationStateSpec
type NodeAllocationStateList = nascrd.NodeAllocationStateList

type NodeAllocationState struct {
	*nascrd.NodeAllocationState
	clientset nvclientset.Interface
}

func NewNodeAllocationState(config *NodeAllocationStateConfig, clientset nvclientset.Interface) *NodeAllocationState {
	object := &nascrd.NodeAllocationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	if config.Owner != nil {
		object.OwnerReferences = []metav1.OwnerReference{*config.Owner}
	}

	nascrd := &NodeAllocationState{
		object,
		clientset,
	}

	return nascrd
}

func (g *NodeAllocationState) GetOrCreate() error {
	err := g.Get()
	if err == nil {
		return nil
	}
	if errors.IsNotFound(err) {
		return g.Create()
	}
	return err
}

func (g *NodeAllocationState) Create() error {
	crd := g.NodeAllocationState.DeepCopy()
	crd, err := g.clientset.NasV1alpha1().NodeAllocationStates(g.Namespace).Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	*g.NodeAllocationState = *crd
	return nil
}

func (g *NodeAllocationState) Delete() error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	err := g.clientset.NasV1alpha1().NodeAllocationStates(g.Namespace).Delete(context.TODO(), g.NodeAllocationState.Name, deleteOptions)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (g *NodeAllocationState) Update(spec *nascrd.NodeAllocationStateSpec) error {
	crd := g.NodeAllocationState.DeepCopy()
	crd.Spec = *spec
	crd, err := g.clientset.NasV1alpha1().NodeAllocationStates(g.Namespace).Update(context.TODO(), crd, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	*g.NodeAllocationState = *crd
	return nil
}

func (g *NodeAllocationState) UpdateStatus(status string) error {
	crd := g.NodeAllocationState.DeepCopy()
	crd.Status = status
	crd, err := g.clientset.NasV1alpha1().NodeAllocationStates(g.Namespace).Update(context.TODO(), crd, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	*g.NodeAllocationState = *crd
	return nil
}

func (g *NodeAllocationState) Get() error {
	crd, err := g.clientset.NasV1alpha1().NodeAllocationStates(g.Namespace).Get(context.TODO(), g.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	*g.NodeAllocationState = *crd
	return nil
}