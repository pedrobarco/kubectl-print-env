package client

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) getSecret(name string) (*v1.Secret, error) {
	s, err := c.Clientset.CoreV1().Secrets(c.Namespace).Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return s, nil
}

func (c *Client) fromSecretRef(sks *v1.SecretEnvSource) []v1.EnvVar {
	return c.FromSecret(sks.Name)
}

func (c *Client) fromSecretKeyRef(sks *v1.SecretKeySelector) v1.EnvVar {
	s, err := c.getSecret(sks.Name)
	if err != nil {
		return v1.EnvVar{}
	}

	return v1.EnvVar{Name: s.Name, Value: string(s.Data[sks.Key])}
}

func (c *Client) FromSecret(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	s, err := c.getSecret(name)
	if err != nil {
		return out
	}

	for k, v := range s.Data {
		out = append(out, v1.EnvVar{Name: k, Value: string(v)})
	}

	return out
}
