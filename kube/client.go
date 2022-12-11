package kube

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nkzren/ecos/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeClient typev1.CoreV1Interface
var kubeConfig config.KubeConf

func init() {
	var err error
	kubeConfig = config.Root.Kube
	kubeClient, err = getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getClient() (typev1.CoreV1Interface, error) {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1(), nil
}

func getConfig() (*restclient.Config, error) {
	switch kubeConfig.Env {
	case "dev":
		return clientcmd.BuildConfigFromFlags("", kubeConfig.ConfPath)
	case "kube":
		return clientcmd.BuildConfigFromFlags("", "")
	default:
		return nil, errors.New("Invalid kube env")
	}
}

func GetNodes() (*corev1.NodeList, error) {
	return kubeClient.Nodes().List(context.Background(), metav1.ListOptions{})
}

func UpdateLabel(node *corev1.Node, value string) error {
	node.ObjectMeta.Labels["ecos"] = value
	_, err := kubeClient.Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	return err
}
