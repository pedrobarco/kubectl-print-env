package parser

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Parser) getPod(name string) (*v1.Pod, error) {
	pod, err := p.Clientset.CoreV1().Pods(p.Namespace).Get(p.Context, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}
	return pod, nil
}

func (p *Parser) FromPod(name string) []v1.EnvVar {
	out := []v1.EnvVar{}

	po, err := p.getPod(name)
	if err != nil {
		return out
	}

	for _, ct := range po.Spec.Containers {
		out = append(out, p.fromEnv(ct.Env, po)...)
		out = append(out, p.fromEnvFrom(ct.EnvFrom)...)
	}

	return out
}
