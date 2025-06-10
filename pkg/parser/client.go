package parser

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type client struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewClient(configFlags *genericclioptions.ConfigFlags, namespace string) (*client, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Parserset: %w", err)
	}

	return &client{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func (c *client) GetDeployment(ctx context.Context, name string) (*appsv1.Deployment, error) {
	deploy, err := c.clientset.AppsV1().Deployments(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deploy, nil
}

func (c *client) GetPod(ctx context.Context, name string) (*v1.Pod, error) {
	pod, err := c.clientset.CoreV1().Pods(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}
	return pod, nil
}

func (c *client) GetSecret(ctx context.Context, name string) (*v1.Secret, error) {
	s, err := c.clientset.CoreV1().Secrets(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	return s, nil
}

func (c *client) GetConfigMap(ctx context.Context, name string) (*v1.ConfigMap, error) {
	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}
	return cm, nil
}

func (c *client) GetJob(ctx context.Context, name string) (*batchv1.Job, error) {
	job, err := c.clientset.BatchV1().Jobs(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}
