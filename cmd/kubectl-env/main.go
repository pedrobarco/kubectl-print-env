package main

import (
	"github.com/pedrobarco/kubectl-env/cmd/kubectl-env/env"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	c := env.NewCmdEnv()
	_ = c.Execute()
}
