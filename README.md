# kubectl-print-env

A kubectl plugin for building config files from k8s environments

With `kubectl-print-env`, it is possible to create a config file from kubernetes
resources. This plugin prints configs by parsing environment information about
the specified resources. You can select the output format using the `--output`
flag.


## Installation

`kubectl-print-env` can be installed using [Krew](https://krew.sigs.k8s.io):

```bash
$ kubectl krew install env
```

or by downloading the binary from the [releases](https://github.com/pedrobarco/kubectl-print-env/releases) page.

Alternatively, `kubectl-print-env` can be installed by running

```bash
$ go install github.com/pedrobarco/kubectl-print-env
```

or by cloning this repository and running:

```bash
$ make build && sudo make install
```


## Usage

```
kubectl print-env [(-o|--output=)dotenv|json|toml|yaml] (TYPE[.VERSION][.GROUP] [NAME] | TYPE[.VERSION][.GROUP]/NAME) [flags]

Examples:
  # Build a dotenv config file from a pod
  kubectl print-env pods my-pod

  # Build a JSON config file from a deployment, in the "v1" version of the "apps" API group
  kubectl print-env deployments.v1.apps my-deployment -o json

  # Build a YAML config file from a configmap
  kubectl print-env cm/my-configmap -o yaml

  # Build a TOML config file from a secret, decoding secret values
  kubectl print-env secret my-secret -o toml
```

## Specification

This plugin supports the following resource types:
- [x] ConfigMap
- [x] Secret
- [x] Pod
- [ ] Daemonset
- [ ] Replicaset
- [ ] Statefulset
- [x] Deployment
- [x] Job
- [ ] CronJob
- [ ] Service
- [ ] Ingress

> NOTE: When running `kubectl-print-env`, only resources of this type will be checked
