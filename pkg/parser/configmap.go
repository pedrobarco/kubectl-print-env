package parser

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

func (p *Parser) fromConfigMapRef(ctx context.Context, cmes *v1.ConfigMapEnvSource) ([]v1.EnvVar, error) {
	return p.FromConfigMap(ctx, cmes.Name)
}

func (p *Parser) fromConfigMapKeyRef(ctx context.Context, cmks *v1.ConfigMapKeySelector) (*v1.EnvVar, error) {
	cm, err := p.client.GetConfigMap(ctx, cmks.Name)
	if err != nil {
		return nil, err
	}
	return &v1.EnvVar{Name: cmks.Key, Value: cm.Data[cmks.Key]}, nil
}

func (p *Parser) FromConfigMap(ctx context.Context, name string) ([]v1.EnvVar, error) {
	cm, err := p.client.GetConfigMap(ctx, name)
	if err != nil {
		return nil, err
	}

	var out []v1.EnvVar
	for k, v := range cm.Data {
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}

	return out, nil
}
