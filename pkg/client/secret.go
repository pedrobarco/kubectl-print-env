package client

import (
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *client) getSecret(name string) (*v1.Secret, error) {
	s, err := c.Clientset.CoreV1().Secrets(c.Namespace).Get(c.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	return s, nil
}

func (c *client) fromSecretRef(sks *v1.SecretEnvSource) []v1.EnvVar {
	// TODO get secret from client package
	return nil
}

func (c *client) fromSecretKeyRef(sks *v1.SecretKeySelector) v1.EnvVar {
	// TODO get secret from client package
	return v1.EnvVar{}
}

func (c *client) FromSecret(name string) []corev1.EnvVar {
	out := []v1.EnvVar{}

	s, err := c.getSecret(name)
	if err != nil {
		return out
	}

	for k, dv := range s.Data {
		v, err := base64.StdEncoding.DecodeString(string(dv))
		if err != nil {
			fmt.Println("failed to decode string: %w", err)
			continue
		}

		out = append(out, corev1.EnvVar{Name: k, Value: string(v)})
	}

	return out
}
