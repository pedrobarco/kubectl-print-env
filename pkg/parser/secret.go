package parser

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Parser) getSecret(name string) (*v1.Secret, error) {
	s, err := p.Clientset.CoreV1().Secrets(p.Namespace).Get(p.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return s, nil
}

func (p *Parser) fromSecretRef(sks *v1.SecretEnvSource) []v1.EnvVar {
	return p.FromSecret(sks.Name)
}

func (p *Parser) fromSecretKeyRef(sks *v1.SecretKeySelector) v1.EnvVar {
	s, err := p.getSecret(sks.Name)
	if err != nil {
		return v1.EnvVar{}
	}

	return v1.EnvVar{Name: s.Name, Value: string(s.Data[sks.Key])}
}

func (p *Parser) FromSecret(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	s, err := p.getSecret(name)
	if err != nil {
		return out
	}

	for k, v := range s.Data {
		out = append(out, v1.EnvVar{Name: k, Value: string(v)})
	}

	return out
}
