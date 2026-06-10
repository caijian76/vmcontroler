package vm

import (
	"context"
	"fmt"
	"log"
	"time"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	watch, err := VirtClient.VirtualMachine(Namespace).Watch(ctx, k8smetav1.ListOptions{
		Watch:         true,
		FieldSelector: "metadata.name=" + vmname,
	})
	if err != nil {
		log.Println("关闭VM:"+vmname+"监听状态失败:", err)
		return err
	}
	defer watch.Stop()

	for event := range watch.ResultChan() {
		if event.Object.(*v1.VirtualMachine).Status.PrintableStatus == "Stopped" {
			log.Println("关闭VM:" + vmname + "成功")
			return nil
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("关闭VM:" + vmname + "超时，未确认虚拟机已停止")
	}

	return nil
}
