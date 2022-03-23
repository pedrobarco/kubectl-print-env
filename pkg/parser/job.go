package parser

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Parser) getJob(name string) (*batchv1.Job, error) {
	job, err := p.Clientset.BatchV1().Jobs(p.Namespace).Get(p.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

func (p *Parser) FromJob(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	j, err := p.getJob(name)
	if err != nil {
		return out
	}

	for _, ct := range j.Spec.Template.Spec.Containers {
		out = append(out, p.fromEnv(ct.Env)...)
		out = append(out, p.fromEnvFrom(ct.EnvFrom)...)
	}

	return out
}
