package printenv

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pedrobarco/kubectl-print-env/pkg/printers"
)

type FormatFlags struct {
	Value  string
	Format printers.Format
}

func (o *FormatFlags) AllowedFormats() []string {
	formats := []string{"dotenv", "json", "yaml", "toml"}
	sort.Strings(formats)
	return formats
}

func (o FormatFlags) String() string {
	return o.Value
}

func (o *FormatFlags) Type() string {
	return "format"
}

func (o *FormatFlags) Set(v string) error {
	v = strings.ToLower(v)
	switch v {
	case "dotenv":
		o.Format = printers.DotEnv
	case "json":
		o.Format = printers.Json
	case "yaml":
		o.Format = printers.Yaml
	case "toml":
		o.Format = printers.Toml
	default:
		return fmt.Errorf("unable to match a printer suitable for the output format %q", v)
	}
	o.Value = v
	return nil
}
