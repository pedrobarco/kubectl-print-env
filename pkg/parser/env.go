package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/util/fieldpath"
)

func (p *Parser) fromEnv(ctx context.Context, env []v1.EnvVar, obj interface{}) ([]v1.EnvVar, error) {
	var out []v1.EnvVar
	for _, ev := range env {
		k := ev.Name
		v := ev.Value
		if evs := ev.ValueFrom; evs != nil {
			from, err := p.fromValueFrom(ctx, evs, obj)
			if err != nil {
				return nil, err
			}
			v = from.Value
		}
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}
	return out, nil
}

func (p *Parser) fromValueFrom(ctx context.Context, evs *v1.EnvVarSource, obj interface{}) (*v1.EnvVar, error) {
	if cmks := evs.ConfigMapKeyRef; cmks != nil {
		return p.fromConfigMapKeyRef(ctx, cmks)
	}
	if sks := evs.SecretKeyRef; sks != nil {
		return p.fromSecretKeyRef(ctx, sks)
	}
	if ofs := evs.FieldRef; ofs != nil {
		return p.fromFieldRef(ctx, ofs, obj)
	}
	return nil, fmt.Errorf("unsupported EnvVarSource type: %T", evs)
}

func (p *Parser) fromFieldRef(ctx context.Context, ofs *v1.ObjectFieldSelector, obj interface{}) (*v1.EnvVar, error) {
	value, err := fieldpath.ExtractFieldPathAsString(obj, ofs.FieldPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract field path %s: %w", ofs.FieldPath, err)
	}
	return &v1.EnvVar{Name: ofs.FieldPath, Value: value}, nil
}

func (p *Parser) fromEnvFrom(ctx context.Context, env []v1.EnvFromSource) ([]v1.EnvVar, error) {
	var out []v1.EnvVar
	for _, efs := range env {
		if cmes := efs.ConfigMapRef; cmes != nil {
			evs, err := p.fromConfigMapRef(ctx, efs.ConfigMapRef)
			if err != nil {
				return nil, fmt.Errorf("failed to get env from configmap ref %s: %w", cmes.Name, err)
			}
			out = append(out, evs...)
			continue
		}

		if ses := efs.SecretRef; ses != nil {
			evs, err := p.fromSecretRef(ctx, efs.SecretRef)
			if err != nil {
				return nil, fmt.Errorf("failed to get env from secret ref %s: %w", ses.Name, err)
			}
			out = append(out, evs...)
			continue
		}

		return nil, fmt.Errorf("unsupported EnvFromSource type: %T", efs)
	}

	return out, nil
}
