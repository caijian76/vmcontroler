package vm

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualMachine struct {
	Name      string
	Run       string
	Ready     bool
	Status    string
	CPU       string
	Memory    string
	StartTime string
	NodeName  string
}

type VirtualMachineList []VirtualMachine

func ListVm() (*VirtualMachineList, error) {
	var vml VirtualMachineList
	// Fetch list of VMs & VMIs
	vmList, err := VirtClient.VirtualMachine(Namespace).List(context.Background(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot obtain KubeVirt vm list: %v", err)
	}

	vmiList, err := VirtClient.VirtualMachineInstance(Namespace).List(context.Background(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot obtain KubeVirt vmi list: %v", err)
	}

	vmiStartTimes := make(map[string]time.Time, len(vmiList.Items))
	vmiNodeNames := make(map[string]string, len(vmiList.Items))

	for _, vmi := range vmiList.Items {
		if !vmi.CreationTimestamp.IsZero() {
			vmiStartTimes[vmi.Name] = vmi.CreationTimestamp.Time
		}
		vmiNodeNames[vmi.Name] = vmi.Status.NodeName
	}

	for _, vm := range vmList.Items {
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

		if vm.Status.PrintableStatus == "Running" {
			if createdAt, ok := vmiStartTimes[vm.Name]; ok && !createdAt.IsZero() {
				cst := time.FixedZone("CST", 8*60*60)
				startTime = createdAt.In(cst).Format("2006-01-02 15:04:05 CST")
			}
		}

		vml = append(vml, VirtualMachine{
			Name:      vm.Name,
			Run:       string(vm.Status.RunStrategy),
			Ready:     vm.Status.Ready,
			Status:    string(vm.Status.PrintableStatus),
			CPU:       cpu + " vCPU",
			Memory:    memory,
			StartTime: startTime,
			NodeName:  vmiNodeNames[vm.Name],
		})
	}
	return &vml, nil

}
