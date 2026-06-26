package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"vmcontroller/utils"

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

func CreateVM(parse *utils.CreateVMRequest) error {
	vm := NewVM(parse)
	vmjson, _ := (json.Marshal(vm))
	log.Println(string(vmjson))
	return createVM(vm)
}

func createVM(vm *v1.VirtualMachine) error {

	_, err := VirtClient.VirtualMachine(DefaultNamespace).Create(context.Background(), vm, k8smetav1.CreateOptions{})
	return err
}

func NewVM(parse *utils.CreateVMRequest) *v1.VirtualMachine {

	vm := &v1.VirtualMachine{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name: parse.VMName,
		},
		Spec: v1.VirtualMachineSpec{
			Running: BoolPtr(parse.AutoStart),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: k8smetav1.ObjectMeta{
					Labels: map[string]string{
						"kubevirt.io/domain": parse.VMName,
					},
				},
				Spec: v1.VirtualMachineInstanceSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": parse.Node,
					},
					Domain: v1.DomainSpec{
						CPU: &v1.CPU{
							Cores: uint32(parse.CPU),
							// Model: "host-passthrough",
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
										Bus: v1.DiskBus("virtio"),
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
								"memory": resource.MustParse(strconv.Itoa(int(parse.Memory)) + "Gi"),
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
									Path:     fmt.Sprintf("/root/kubevirt/template/%s/%s.img", parse.VMName, parse.VMName),
									Type:     v1.HostDiskType("DiskOrCreate"),
									Capacity: resource.MustParse(strconv.Itoa(int(parse.DiskSize)) + "Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	if parse.MountISO {
		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, v1.Disk{
			Name:      "cdromiso",
			BootOrder: UintPtr(uint(2)),
			DiskDevice: v1.DiskDevice{
				CDRom: &v1.CDRomTarget{
					Bus: v1.DiskBus("sata"),
				},
			},
		})
		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, v1.Volume{
			Name: "cdromiso",
			VolumeSource: v1.VolumeSource{
				HostDisk: &v1.HostDisk{
					Path: parse.ISOPath,
					Type: v1.HostDiskType("Disk"),
				},
			},
		})
	}

	return vm
}
