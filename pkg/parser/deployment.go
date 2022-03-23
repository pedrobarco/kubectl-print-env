package parser

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Parser) getDeployment(name string) (*appsv1.Deployment, error) {
	deploy, err := p.Clientset.AppsV1().Deployments(p.Namespace).Get(p.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deploy, nil
}

func (p *Parser) FromDeployment(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	d, err := p.getDeployment(name)
	if err != nil {
		return out
	}

	for _, ct := range d.Spec.Template.Spec.Containers {
		out = append(out, p.fromEnv(ct.Env)...)
		out = append(out, p.fromEnvFrom(ct.EnvFrom)...)
	}

	return out
}
