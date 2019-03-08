package cmd

import (
	"os"

	"github.com/dcos/dcos-cli/api"
	"github.com/spf13/cobra"
)

const annotationUsageOptions string = "usage_options"

// NewMesosCommand creates the `dcos` command with all the available subcommands.
func NewMesosCommand(ctx api.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "mesos",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	cmd.AddCommand(
		newCmdMesosSandbox(ctx),
	)

	// From the DC/OS CLI docs:
	// When a plugin is called it will have the top level command name as its
	// second argument regardless of whether the plugin only has one executable.
	args := os.Args
	if len(args) > 1 && args[1] == "mesos" {
		if len(args) == 2 {
			cmd.SetArgs([]string{""})
		} else {
			cmd.SetArgs(args[2:])
		}
	}

	// This follows the CLI design guidelines for help formatting.
	cmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{.Name}}
      {{.Short}}{{end}}{{end}}{{end}}{{if or .HasAvailableLocalFlags (ne (index .Annotations "` + annotationUsageOptions + `") "")}}

Options:{{if ne (index .Annotations "` + annotationUsageOptions + `") ""}}{{index .Annotations "` + annotationUsageOptions + `"}}{{else}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	return cmd
}
