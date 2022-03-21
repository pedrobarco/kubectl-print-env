package client

import v1 "k8s.io/api/core/v1"

func (c *Client) fromEnv(env []v1.EnvVar) []v1.EnvVar {
	out := []v1.EnvVar{}
	for _, ev := range env {
		k := ev.Name
		var v string
		if evs := ev.ValueFrom; evs != nil {
			v = c.fromValueFrom(evs).Value
		} else {
			v = ev.Value
		}
		out = append(out, v1.EnvVar{Name: k, Value: v})
	}
	return out
}

func (c *Client) fromValueFrom(evs *v1.EnvVarSource) v1.EnvVar {
	if cmks := evs.ConfigMapKeyRef; cmks != nil {
		return c.fromConfigMapKeyRef(cmks)
	}
	if sks := evs.SecretKeyRef; sks != nil {
		return c.fromSecretKeyRef(sks)
	}
	return v1.EnvVar{}
}

func (c *Client) fromEnvFrom(env []v1.EnvFromSource) []v1.EnvVar {
	out := []v1.EnvVar{}

	for _, efs := range env {
		if cmes := efs.ConfigMapRef; cmes != nil {
			out = append(out, c.fromConfigMapRef(cmes)...)
		}
		if ses := efs.SecretRef; ses != nil {
			out = append(out, c.fromSecretRef(ses)...)
		}
	}

	return out
}
