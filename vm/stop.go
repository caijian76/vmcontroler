package vm

import (
	"context"
	"log"
	"time"

	v1 "kubevirt.io/api/core/v1"
)

func StopVM(vmname string, delay time.Duration) error {
	time.Sleep(delay)
	log.Println("关闭VM:" + vmname + "中...")
	err := VirtClient.VirtualMachine(Namespace).Stop(context.Background(), vmname, &v1.StopOptions{})
	if err != nil {
		log.Println("关闭VM:" + vmname + "失败! " + err.Error())
		return err
	}
	log.Println("关闭VM:" + vmname + "成功")
	return nil
}
