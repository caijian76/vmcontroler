package vm

import (
	"context"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StatusVM(vmname string) (*VirtualMachine, error) {
	repeatvm := VirtualMachine{}
	vm, err := VirtClient.VirtualMachine(Namespace).Get(context.Background(), vmname, k8smetav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	vmi, err := VirtClient.VirtualMachineInstance(Namespace).Get(context.Background(), vmname, k8smetav1.GetOptions{})
	cpu := "0"
	memory := ""
	startTime := ""

	if vm.Spec.Template != nil && vm.Spec.Template.Spec.Domain.CPU != nil {
		cpu = strconv.Itoa(int(vm.Spec.Template.Spec.Domain.CPU.Cores))
	}

	if vm.Spec.Template != nil {
		if q, ok := vm.Spec.Template.Spec.Domain.Resources.Requests[corev1.ResourceMemory]; ok && !q.IsZero() {
			memory = q.String()
		} else if q, ok := vm.Spec.Template.Spec.Domain.Resources.Limits[corev1.ResourceMemory]; ok && !q.IsZero() {
			memory = q.String()
		}
	}

	if err == nil && vm.Status.PrintableStatus == "Running" && !vmi.CreationTimestamp.IsZero() {
		cst := time.FixedZone("CST", 8*60*60)
		startTime = vmi.CreationTimestamp.In(cst).Format("2006-01-02 15:04:05 CST")
	}

	repeatvm.Name = vm.Name
	repeatvm.Run = string(vm.Status.RunStrategy)
	repeatvm.Status = string(vm.Status.PrintableStatus)
	repeatvm.Ready = vm.Status.Ready
	repeatvm.CPU = cpu + " vCPU"
	repeatvm.Memory = memory
	repeatvm.StartTime = startTime

	return &repeatvm, nil

}
