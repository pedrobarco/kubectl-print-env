package parser

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/util/fieldpath"
)

func (p *Parser) fromEnv(env []v1.EnvVar, obj interface{}) []v1.EnvVar {
	out := []v1.EnvVar{}
	for _, ev := range env {
		k := ev.Name
		var v string
		if evs := ev.ValueFrom; evs != nil {
			v = p.fromValueFrom(evs, obj).Value
		} else {
			v = ev.Value
		}
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}
	return out
}

func (p *Parser) fromValueFrom(evs *v1.EnvVarSource, obj interface{}) v1.EnvVar {
	if cmks := evs.ConfigMapKeyRef; cmks != nil {
		return p.fromConfigMapKeyRef(cmks)
	}
	if sks := evs.SecretKeyRef; sks != nil {
		return p.fromSecretKeyRef(sks)
	}
	if ofs := evs.FieldRef; ofs != nil {
		return p.fromFieldRef(ofs, obj)
	}
	return v1.EnvVar{}
}

func (p *Parser) fromFieldRef(ofs *v1.ObjectFieldSelector, obj interface{}) v1.EnvVar {
	value, err := fieldpath.ExtractFieldPathAsString(obj, ofs.FieldPath)
	if err != nil {
		return v1.EnvVar{}
	}

	return v1.EnvVar{Name: ofs.FieldPath, Value: value}
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
