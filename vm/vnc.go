package vm

import v1 "kubevirt.io/client-go/kubevirt/typed/core/v1"

func Getvnc(vmname string) (v1.StreamInterface, error) {
	return VirtClient.VirtualMachineInstance(Namespace).VNC(vmname)
}
