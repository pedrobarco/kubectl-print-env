package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func (p *Parser) FromDeployment(ctx context.Context, name string) ([]v1.EnvVar, error) {
	d, err := p.client.GetDeployment(ctx, name)
	if err != nil {
		return nil, err
	}

	var out []v1.EnvVar
	for _, ct := range d.Spec.Template.Spec.Containers {
		evs, err := p.fromEnv(ctx, ct.Env, d)
		if err != nil {
			return nil, fmt.Errorf("failed to parse env vars from container %s: %w", ct.Name, err)
		}
		out = append(out, evs...)

		evs, err = p.fromEnvFrom(ctx, ct.EnvFrom)
		if err != nil {
			return nil, fmt.Errorf("failed to parse env from in container %s: %w", ct.Name, err)
		}
		out = append(out, evs...)
	}

	return out, nil
}
