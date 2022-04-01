package parser

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type Parser struct {
	ConfigFlags *genericclioptions.ConfigFlags
	Clientset   *kubernetes.Clientset
	Context     context.Context
	Namespace   string
}

func CreateParser(ns string, configFlags *genericclioptions.ConfigFlags) (*Parser, error) {

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Parserset: %w", err)
	}

	return &Parser{
		ConfigFlags: configFlags,
		Clientset:   clientset,
		Context:     context.Background(),
		Namespace:   ns,
	}, nil
}
