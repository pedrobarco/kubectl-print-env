package main

import (
	"github.com/pedrobarco/kubectl-print-env/cmd/kubectl-print-env/printenv"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	c := printenv.NewCmdEnv()
	_ = c.Execute()
}
