package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func (p *Parser) FromPod(ctx context.Context, name string) ([]v1.EnvVar, error) {
	pod, err := p.client.GetPod(ctx, name)
	if err != nil {
		return nil, err
	}

	var out []v1.EnvVar
	for _, ct := range pod.Spec.Containers {
		evs, err := p.fromEnv(ctx, ct.Env, pod)
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
