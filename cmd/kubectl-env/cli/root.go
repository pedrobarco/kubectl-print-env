package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/pedrobarco/kubectl-env/pkg/parser"
	"github.com/pedrobarco/kubectl-env/pkg/printer"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

type Options struct {
	namespace string
	args      []string
	builder   *resource.Builder
	flags     *genericclioptions.ConfigFlags
	out       io.Writer
	parser    *parser.Parser
}

func CheckErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func NewCmdEnv() *cobra.Command {
	o := Options{out: os.Stdout}
	f := genericclioptions.NewConfigFlags(true)
	cmd := &cobra.Command{
		Use:           "kubectl env TYPE[.VERSION][.GROUP] NAME",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(o.Complete(f, cmd, args))
			CheckErr(o.Validate())
			CheckErr(o.Run())
		},
	}
	flags := cmd.Flags()
	f.AddFlags(flags)
	return cmd
}

func (o *Options) Complete(f *genericclioptions.ConfigFlags, cmd *cobra.Command, args []string) error {
	ns, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	p, err := parser.CreateParser(f)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	o.namespace = ns
	o.args = args
	o.flags = f
	o.builder = resource.NewBuilder(f)
	o.parser = p
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run() error {
	result := o.builder.Unstructured().
		NamespaceParam(o.namespace).
		DefaultNamespace().
		ResourceTypeOrNameArgs(true, o.args...).
		Latest().
		Do()

	if err := result.Err(); err != nil {
		return err
	}

	return result.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		var env []v1.EnvVar
		switch info.Mapping.GroupVersionKind.Kind {
		case "Deployment":
			env = o.parser.FromDeployment(info.Name)
		case "ConfigMap":
			env = o.parser.FromConfigMap(info.Name)
		case "Secret":
			env = o.parser.FromSecret(info.Name)
		case "Job":
			env = o.parser.FromJob(info.Name)
		case "Pod":
			env = o.parser.FromPod(info.Name)
		default:
			env = nil
		}

		fmt.Println(printer.Print(env, printer.DotEnv))
		return nil
	})
}
