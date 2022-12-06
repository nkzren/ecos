package kube

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nkzren/ecoscheduler/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var client typev1.CoreV1Interface

func init() {
	var err error
	client, err = getClient(config.Root.Kube.ConfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getClient(cfgPath string) (typev1.CoreV1Interface, error) {
	kubeconfig := filepath.Clean(cfgPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1(), nil
}

func GetNodes() (*corev1.NodeList, error) {
	client, err := getClient(config.Root.Kube.ConfPath)
	if err != nil {
		return nil, err
	}
	return client.Nodes().List(context.Background(), metav1.ListOptions{})
}

func UpdateLabel(node *corev1.Node, value string) error {
	node.ObjectMeta.Labels["ecos"] = value
	client.Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	return nil
}
