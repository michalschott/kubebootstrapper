package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/michalschott/kubebootstrapper/pkg/version"
	"github.com/spf13/cobra"
)

const defaultVersionString = "unavailable"

func newCmdVersion() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			configureAndRunVersion(cmd.Context(), os.Stdout)
		},
	}

	return cmd
}

func configureAndRunVersion(
	ctx context.Context,
	stdout io.Writer,
) {
	fmt.Fprintf(stdout, "Version: %s\n", version.Version)
}
