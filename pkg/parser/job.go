package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func (p *Parser) FromJob(ctx context.Context, name string) ([]v1.EnvVar, error) {
	j, err := p.client.GetJob(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get job %s: %w", name, err)
	}

	var out []v1.EnvVar
	for _, ct := range j.Spec.Template.Spec.Containers {
		evs, err := p.fromEnv(ctx, ct.Env, j)
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
