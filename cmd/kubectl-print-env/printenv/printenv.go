package printenv

import (
	"fmt"
	"os"
	"strings"

	"github.com/pedrobarco/kubectl-print-env/pkg/parser"
	"github.com/pedrobarco/kubectl-print-env/pkg/printers"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

type PrintEnvOptions struct {
	formatFlags *FormatFlags
	configFlags *genericclioptions.ConfigFlags

	args      []string
	namespace string

	builder *resource.Builder
	parser  *parser.Parser
}

var (
	envLong = templates.LongDesc(i18n.T(`
		Build config files from k8s environments.

		Prints a config file by parsing environment information about the
		specified resources.
		You can select the output format using the --output flag.`))

	envExample = templates.Examples(i18n.T(`
		# Build a dotenv config file from a pod
		kubectl print-env pods my-pod

		# Build a JSON config file from a deployment, in the "v1" version of the "apps" API group
		kubectl print-env deployments.v1.apps my-deployment -o json

		# Build a YAML config file from a configmap
		kubectl print-env cm/my-configmap -o yaml

		# Build a TOML config file from a secret, decoding secret values
		kubectl print-env secret my-secret -o toml`))
)

func CheckErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func NewCmdEnv() *cobra.Command {
	o := &PrintEnvOptions{formatFlags: &FormatFlags{Format: printers.DotEnv}}
	f := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:          fmt.Sprintf("kubectl print-env [(-o|--output=)%s] (TYPE[.VERSION][.GROUP] [NAME] | TYPE[.VERSION][.GROUP]/NAME)", strings.Join(o.formatFlags.AllowedFormats(), "|")),
		Short:        i18n.T("Build config files from k8s environments"),
		Long:         envLong,
		Example:      envExample,
		SilenceUsage: true,
		Args:         cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(o.Complete(f, cmd, args))
			CheckErr(o.Validate())
			CheckErr(o.Run())
		},
	}

	flags := cmd.Flags()
	f.AddFlags(flags)
	cmd.Flags().VarP(o.formatFlags, "output", "o", fmt.Sprintf("Output format. One of: %s", strings.Join(o.formatFlags.AllowedFormats(), "|")))
	return cmd
}

func (o *PrintEnvOptions) Complete(f *genericclioptions.ConfigFlags, cmd *cobra.Command, args []string) error {
	ns, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	p, err := parser.CreateParser(f)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	o.namespace = ns
	o.parser = p
	o.args = args
	o.configFlags = f
	o.builder = resource.NewBuilder(f)
	return nil
}

func (o *PrintEnvOptions) Validate() error {
	return nil
}

func (o *PrintEnvOptions) Run() error {
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

		fmt.Println(printers.Print(env, o.formatFlags.Format))
		return nil
	})
}
