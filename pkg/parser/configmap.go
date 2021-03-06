package parser

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Parser) getConfigMap(name string) (*v1.ConfigMap, error) {
	cm, err := p.Clientset.CoreV1().ConfigMaps(p.Namespace).Get(p.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}
	return cm, nil
}

func (p *Parser) fromConfigMapRef(cmes *v1.ConfigMapEnvSource) []v1.EnvVar {
	return p.FromConfigMap(cmes.Name)
}

func (p *Parser) fromConfigMapKeyRef(cmks *v1.ConfigMapKeySelector) v1.EnvVar {
	cm, err := p.getConfigMap(cmks.Name)
	if err != nil {
		return v1.EnvVar{}
	}
	return v1.EnvVar{Name: cmks.Key, Value: cm.Data[cmks.Key]}
}

func (p *Parser) FromConfigMap(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	cm, err := p.getConfigMap(name)
	if err != nil {
		return out
	}

	for k, v := range cm.Data {
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}

	return out
}
