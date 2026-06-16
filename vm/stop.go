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
	log.Printf("关闭VM:%s中...", vmname)
	time.Sleep(delay)
	err := VirtClient.VirtualMachine(Namespace).Stop(context.Background(), vmname, &v1.StopOptions{})
	if err != nil {
		log.Printf("关闭VM:%s失败! %v", vmname, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	watch, err := VirtClient.VirtualMachine(Namespace).Watch(ctx, k8smetav1.ListOptions{
		Watch:         true,
		FieldSelector: fmt.Sprintf("metadata.name=%s", vmname),
	})
	if err != nil {
		log.Printf("关闭VM:%s监听状态失败:%v", vmname, err)
		return err
	}
	defer watch.Stop()

	for event := range watch.ResultChan() {
		if event.Object.(*v1.VirtualMachine).Status.PrintableStatus == "Stopped" {
			log.Printf("关闭VM:%s成功", vmname)
			return nil
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("关闭VM:%s超时，未确认虚拟机已停止", vmname)
	}

	return nil
}
