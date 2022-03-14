package client

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *client) getConfigMap(name string) (*corev1.ConfigMap, error) {
	cm, err := c.Clientset.CoreV1().ConfigMaps(c.Namespace).Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}
	return cm, nil
}

func (c *client) fromConfigMapRef(sks *v1.ConfigMapEnvSource) []v1.EnvVar {
	return nil
}

func (c *client) fromConfigMapKeyRef(cmks *v1.ConfigMapKeySelector) v1.EnvVar {
	return v1.EnvVar{}
}

func (c *client) FromConfigMap(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	cm, err := c.getConfigMap(name)
	if err != nil {
		return out
	}

	for k, v := range cm.Data {
		out = append(out, corev1.EnvVar{Name: k, Value: v})
	}

	return out
}
