package main

import (
	"context"
	"fmt"
	"log"
	"vmcontroller/web"

	"time"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

func main() {

	//vm.WatchVM(virtClient, namespace)
	web.WebStart()
	//vnc(virtClient, namespace, "cc")
}

func CreateVM(virtClient kubecli.KubevirtClient, namespace, vmname string) {
	err := virtClient.VirtualMachine(namespace).Start(context.Background(), vmname, &v1.StartOptions{})
	if err != nil {
		log.Printf("cannot obtain KubeVirt vm list: %v\n", err)
	}

	for {
		checkvm, err := virtClient.VirtualMachine(namespace).Get(context.Background(), vmname, k8smetav1.GetOptions{})
		if err != nil {
			log.Printf("cannot obtain KubeVirt vm list: %v\n", err)
		}
		fmt.Println(checkvm.Status.Conditions[0].Status)
		time.Sleep(5 * time.Second)
		if checkvm.Status.PrintableStatus == "Running" {
			break
		}
	}
}
