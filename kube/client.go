package kube

import (
	"context"
	"log"
	"path/filepath"

	"github.com/nkzren/ecoscheduler/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

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

func GetNodes(c config.KubeConf) (*corev1.NodeList, error) {
	client, err := getClient(c.ConfPath)
	if err != nil {
		return nil, err
	}
	return client.Nodes().List(context.Background(), metav1.ListOptions{})
}
