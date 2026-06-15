package vm

import (
	"context"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNode() ([]string, error) {
	nodes, err := VirtClient.CoreV1().Nodes().List(context.Background(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeNames := make([]string, len(nodes.Items))
	for i, node := range nodes.Items {
		nodeNames[i] = node.Name
	}

	return nodeNames, nil
}
