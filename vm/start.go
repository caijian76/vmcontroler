package vm

import (
	"context"
	"fmt"
	"log"
	"time"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

func StartVM(vmname string) error {
	log.Println("启动VM:" + vmname + "中...")
	err := VirtClient.VirtualMachine(Namespace).Start(context.Background(), vmname, &v1.StartOptions{})
	if err != nil {
		log.Println("启动VM:"+vmname+"失败 ,", err.Error())
		return err
	}

	ctx, canel := context.WithTimeout(context.Background(), 60*time.Second)
	defer canel()
	watch, err := VirtClient.VirtualMachine(Namespace).Watch(ctx, k8smetav1.ListOptions{
		Watch:         true,
		FieldSelector: "metadata.name=" + vmname,
	})
	defer watch.Stop()
	if err != nil {
		log.Println(err)
		return err
	}

	for event := range watch.ResultChan() {
		if event.Object.(*v1.VirtualMachine).Status.Ready {
			break
		}
	}
	if ctx.Err() == context.DeadlineExceeded {

		log.Println("启动VM:" + vmname + "已超时,自动关闭VM")

		go StopVM(vmname, 0)
		return fmt.Errorf("启动VM:" + vmname + "已超时,自动关闭VM")
	}

	log.Println("启动VM:" + vmname + "启动成功")
	return nil
}
