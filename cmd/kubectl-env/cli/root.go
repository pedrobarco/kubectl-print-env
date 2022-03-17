package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/pedrobarco/kubectl-env/pkg/client"
	"github.com/pedrobarco/kubectl-env/pkg/printer"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

type Options struct {
	namespace string
	args      []string
	builder   *resource.Builder
	flags     *genericclioptions.ConfigFlags
	out       io.Writer
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

	o.namespace = ns
	o.args = args
	o.flags = f
	o.builder = resource.NewBuilder(f)
	return nil
}

func (o *Options) Run() error {
	c, err := client.CreateClient(o.flags)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	env := c.FromDeployment(o.args[0])
	fmt.Println(printer.Print(env, printer.DotEnv))

	return nil
}
