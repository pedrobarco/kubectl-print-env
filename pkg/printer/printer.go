package printer

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Format int

const (
	DotEnv Format = iota + 1
	Yaml
	Json
	Toml
)

var printfn = map[Format]func(v1.EnvVar) string{
	DotEnv: toDotEnv,
	Yaml:   toYaml,
	Json:   toJson,
	Toml:   toToml,
}

func Print(env []v1.EnvVar, f Format) string {
	out := strings.Builder{}
	for i, ev := range env {
		out.WriteString(printfn[f](ev))
		if i != len(env)-1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}

func toDotEnv(ev v1.EnvVar) string {
	return fmt.Sprintf("%s=%s", ev.Name, ev.Value)
}

func toYaml(ev v1.EnvVar) string {
	return fmt.Sprintf("%s: \"%s\"", ev.Name, ev.Value)
}

func toJson(ev v1.EnvVar) string {
	return fmt.Sprintf("\"%s\": \"%s\",", ev.Name, ev.Value)
}

func toToml(ev v1.EnvVar) string {
	return fmt.Sprintf("\"%s\" = \"%s\"", ev.Name, ev.Value)
}
