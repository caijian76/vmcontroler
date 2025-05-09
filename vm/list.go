package vm

import (
	"context"
	"fmt"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualMachine struct {
	Name   string
	Run    string
	Ready  bool
	Status string
}

type VirtualMachineList []VirtualMachine

func ListVm() (*VirtualMachineList, error) {
	var vml VirtualMachineList
	// Fetch list of VMs & VMIs
	vmList, err := VirtClient.VirtualMachine(Namespace).List(context.Background(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot obtain KubeVirt vm list: %v", err)
	}

	for _, vm := range vmList.Items {

		vml = append(vml, VirtualMachine{
			Name:   vm.Name,
			Run:    string(vm.Status.RunStrategy),
			Ready:  vm.Status.Ready,
			Status: string(vm.Status.PrintableStatus),
		})

	}
	return &vml, nil

}
