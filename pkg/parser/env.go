package parser

import v1 "k8s.io/api/core/v1"

func (p *Parser) fromEnv(env []v1.EnvVar) []v1.EnvVar {
	out := []v1.EnvVar{}
	for _, ev := range env {
		k := ev.Name
		var v string
		if evs := ev.ValueFrom; evs != nil {
			v = p.fromValueFrom(evs).Value
		} else {
			v = ev.Value
		}
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}
	return out
}

func (p *Parser) fromValueFrom(evs *v1.EnvVarSource) v1.EnvVar {
	if cmks := evs.ConfigMapKeyRef; cmks != nil {
		return p.fromConfigMapKeyRef(cmks)
	}
	if sks := evs.SecretKeyRef; sks != nil {
		return p.fromSecretKeyRef(sks)
	}
	return v1.EnvVar{}
}

func (p *Parser) fromEnvFrom(env []v1.EnvFromSource) []v1.EnvVar {
	out := []v1.EnvVar{}

	for _, efs := range env {
		if cmes := efs.ConfigMapRef; cmes != nil {
			out = append(out, p.fromConfigMapRef(cmes)...)
		}
		if ses := efs.SecretRef; ses != nil {
			out = append(out, p.fromSecretRef(ses)...)
		}
	}

	return out
}
