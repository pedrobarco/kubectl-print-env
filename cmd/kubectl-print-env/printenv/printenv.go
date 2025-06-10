package printenv

import (
	"fmt"
	"strings"

	"github.com/pedrobarco/kubectl-print-env/pkg/parser"
	"github.com/pedrobarco/kubectl-print-env/pkg/printers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

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

type PrintEnvOptions struct {
	formatFlags *FormatFlags
	configFlags *genericclioptions.ConfigFlags

	args      []string
	namespace string

	parser *parser.Parser
}

func NewPrintEnvOptions() *PrintEnvOptions {
	return &PrintEnvOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		formatFlags: &FormatFlags{Format: printers.DotEnv},
	}
}

func NewCmdEnv() *cobra.Command {
	o := NewPrintEnvOptions()

	cmd := &cobra.Command{
		Use:          i18n.T(fmt.Sprintf("kubectl print-env [(-o|--output=)%s] (TYPE[.VERSION][.GROUP] [NAME] | TYPE[.VERSION][.GROUP]/NAME)", strings.Join(o.formatFlags.AllowedFormats(), "|"))),
		Short:        i18n.T("Build config files from k8s environments"),
		Long:         envLong,
		Example:      envExample,
		SilenceUsage: true,
		Args:         cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().VarP(o.formatFlags, "output", "o", fmt.Sprintf("Output format. One of: %s", strings.Join(o.formatFlags.AllowedFormats(), "|")))
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *PrintEnvOptions) Complete(cmd *cobra.Command, args []string) error {
	ns, _, err := o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	parser, err := parser.New(o.configFlags, ns)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	o.args = args
	o.namespace = ns
	o.parser = parser
	return nil
}

func (o *PrintEnvOptions) Validate() error {
	return nil
}

func (o *PrintEnvOptions) Run() error {
	result := resource.NewBuilder(o.configFlags).Unstructured().NamespaceParam(o.namespace).
		ResourceTypeOrNameArgs(false, o.args...).Latest().Do()
	if err := result.Err(); err != nil {
		return err
	}

	return result.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		env, err := o.parser.Parse(info.Mapping.GroupVersionKind.Kind, info.Name)
		if err != nil {
			return fmt.Errorf("error parsing %s %s: %w", info.Mapping.GroupVersionKind.Kind, info.Name, err)
		}

		fmt.Println(printers.Print(env, o.formatFlags.Format))
		return nil
	})
}
