package vm

import (
	"context"
	"encoding/json"
	"log"

	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

func BoolPtr(b bool) *bool {
	return &b
}
func UintPtr(i uint) *uint {
	return &i
}

func CreateVM(vmname string) error {
	vm := NewVM(vmname)
	vmjson, _ := (json.Marshal(vm))
	log.Println(string(vmjson))
	return createVM(vm)
}

func createVM(vm *v1.VirtualMachine) error {

	_, err := VirtClient.VirtualMachine(Namespace).Create(context.Background(), vm, k8smetav1.CreateOptions{})
	return err
}

func NewVM(vmName string) *v1.VirtualMachine {

	return &v1.VirtualMachine{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name: vmName,
		},
		Spec: v1.VirtualMachineSpec{
			Running: BoolPtr(true),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: k8smetav1.ObjectMeta{
					Labels: map[string]string{
						"kubevirt.io/domain": vmName,
					},
				},
				Spec: v1.VirtualMachineInstanceSpec{
					// 节点固定调度到k8s2
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": "k8s2",
					},
					Domain: v1.DomainSpec{
						CPU: &v1.CPU{
							Cores: 2,
							Model: "host-passthrough",
						},
						Devices: v1.Devices{
							Inputs: []v1.Input{
								{
									Type: "tablet",
									Bus:  "usb",
									Name: "tablet1",
								},
							},
							Disks: []v1.Disk{
								{
									Name:      "host-disk",
									BootOrder: UintPtr(uint(1)),
									DiskDevice: v1.DiskDevice{Disk: &v1.DiskTarget{
										Bus: "virtio",
									},
									},
								},
							},
							Interfaces: []v1.Interface{
								{
									Name:  "default",
									Model: "virtio",
									InterfaceBindingMethod: v1.InterfaceBindingMethod{
										Masquerade: &v1.InterfaceMasquerade{},
									},
								},
							},
						},
						Features: &v1.Features{
							ACPI: v1.FeatureState{},
							APIC: &v1.FeatureAPIC{},
							SMM:  &v1.FeatureState{},
						},
						Resources: v1.ResourceRequirements{
							Requests: corev1.ResourceList{
								"memory": resource.MustParse("4Gi"),
							},
						},
					},
					Networks: []v1.Network{
						{
							Name: "default",
							NetworkSource: v1.NetworkSource{
								Pod: &v1.PodNetwork{},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "host-disk",
							VolumeSource: v1.VolumeSource{
								HostDisk: &v1.HostDisk{
									Path:     "/root/kubevirt/template/ubuntu-gui-base-2204/ubuntu-gui-base-2204.img",
									Type:     v1.HostDiskType("DiskOrCreate"),
									Capacity: resource.MustParse("50Gi"),
								},
							},
						},
					},
				},
			},
		},
	}
}
