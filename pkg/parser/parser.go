package parser

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type parserFn func(ctx context.Context, name string) ([]v1.EnvVar, error)

type Parser struct {
	client   *client
	registry map[string]parserFn
}

func New(configFlags *genericclioptions.ConfigFlags, namespace string) (*Parser, error) {
	client, err := NewClient(configFlags, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser client: %w", err)
	}

	parser := &Parser{
		client:   client,
		registry: make(map[string]parserFn),
	}

	parser.registry["Deployment"] = parser.FromDeployment
	parser.registry["Pod"] = parser.FromPod
	parser.registry["Job"] = parser.FromJob
	parser.registry["ConfigMap"] = parser.FromConfigMap
	parser.registry["Secret"] = parser.FromSecret

	return parser, nil
}

func (p *Parser) Parse(kind, name string) ([]v1.EnvVar, error) {
	if fn, ok := p.registry[kind]; ok {
		return fn(context.Background(), name)
	}
	return nil, fmt.Errorf("parser not implemnented for kind %s", kind)
}
