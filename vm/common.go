package vm

import (
	"log"

	"github.com/spf13/pflag"
	"kubevirt.io/client-go/kubecli"
)

var Namespace string
var VirtClient kubecli.KubevirtClient

func init() {
	// kubecli.DefaultClientConfig() prepares config using kubeconfig.
	// typically, you need to set env variable, KUBECONFIG=<path-to-kubeconfig>/.kubeconfig
	var clientConfig = kubecli.DefaultClientConfig(&pflag.FlagSet{})
	var err error
	// retrive default namespace.
	Namespace, _, err = clientConfig.Namespace()
	if err != nil {
		log.Fatalf("error in namespace : %v\n", err)
	}

	// // get the kubevirt client, using which kubevirt resources can be managed.
	VirtClient, err = kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

}
