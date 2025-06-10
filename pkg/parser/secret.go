package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func (p *Parser) fromSecretRef(ctx context.Context, sks *v1.SecretEnvSource) ([]v1.EnvVar, error) {
	return p.FromSecret(ctx, sks.Name)
}

func (p *Parser) fromSecretKeyRef(ctx context.Context, sks *v1.SecretKeySelector) (*v1.EnvVar, error) {
	s, err := p.client.GetSecret(ctx, sks.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", sks.Name, err)
	}

	return &v1.EnvVar{Name: s.Name, Value: string(s.Data[sks.Key])}, nil
}

func (p *Parser) FromSecret(ctx context.Context, name string) ([]v1.EnvVar, error) {
	s, err := p.client.GetSecret(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", name, err)
	}

	var out []v1.EnvVar
	for k, v := range s.Data {
		out = append(out, v1.EnvVar{Name: k, Value: string(v)})
	}

	return out, nil
}
