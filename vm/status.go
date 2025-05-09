package vm

import (
	"context"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StatusVM(vmname string) (*VirtualMachine, error) {
	repeatvm := VirtualMachine{}
	vm, err := VirtClient.VirtualMachine(Namespace).Get(context.Background(), vmname, k8smetav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	repeatvm.Name = vm.Name
	repeatvm.Run = string(vm.Status.RunStrategy)
	repeatvm.Status = string(vm.Status.PrintableStatus)
	repeatvm.Ready = vm.Status.Ready

	return &repeatvm, nil

}
