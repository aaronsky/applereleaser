package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Execute is the primary function to initiate the command line interface for applereleaser
func Execute(version string, exit func(int), args []string) {
	fmt.Println()
	defer fmt.Println()
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd   *cobra.Command
	debug bool
	exit  func(int)
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	var root = &rootCmd{
		exit: exit,
	}

	var cmd = &cobra.Command{
		Use:           "applereleaser",
		Short:         "Manage the release of your Apple App Store applications with ease",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")
	cmd.AddCommand(
		newReleaseCmd().cmd,
		newCheckCmd().cmd,
		newInitCmd().cmd,
	)

	root.cmd = cmd
	return root
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		var code = 1
		var msg = "command failed"
		if eerr, ok := err.(*exitError); ok {
			code = eerr.code
			if eerr.details != "" {
				msg = eerr.details
			}
		}
		fmt.Printf("%s: %s\n", msg, err)
		cmd.exit(code)
	}
}