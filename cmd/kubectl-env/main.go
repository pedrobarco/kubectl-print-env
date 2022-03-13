package main

import (
	"github.com/pedrobarco/kubectl-env/cmd/kubectl-env/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cli.InitAndExecute()
}
