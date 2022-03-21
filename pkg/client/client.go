package client

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	ConfigFlags *genericclioptions.ConfigFlags
	Clientset   *kubernetes.Clientset
	Context     context.Context
	Namespace   string
}

func CreateClient(configFlags *genericclioptions.ConfigFlags) (*Client, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	ns := *configFlags.Namespace
	if ns == "" {
		ns = metav1.NamespaceDefault
	}

	return &Client{ConfigFlags: configFlags, Clientset: clientset, Context: context.Background(), Namespace: ns}, nil
}
