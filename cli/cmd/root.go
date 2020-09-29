package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultKubeBootstrapperNamespace = "kubebootstrapper"
)

var (
	stdout = color.Output
	stderr = color.Error

	okStatus   = color.New(color.FgGreen, color.Bold).SprintFunc()("\u221A")  // √
	warnStatus = color.New(color.FgYellow, color.Bold).SprintFunc()("\u203C") // ‼
	failStatus = color.New(color.FgRed, color.Bold).SprintFunc()("\u00D7")    // ×

	// These regexs are not as strict as they could be, but are a quick and dirty sanity check against illegal characters.
	alphaNumDash = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	kubeBootstrapperNamespace string
	verbose                   bool
)

var RootCmd = &cobra.Command{
	Use:   "kubebootstrapper",
	Short: "kubebootstrapper cli",
	Long: `kubebootstrapper 
			will help you
			keep cluster with good shape`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.PanicLevel)
		}

		kubeBootstrapperNamespaceFromEnv := os.Getenv("KUBEBOOTSTRAPPER_NAMESPACE")
		if kubeBootstrapperNamespace == defaultKubeBootstrapperNamespace && kubeBootstrapperNamespaceFromEnv != "" {
			kubeBootstrapperNamespace = kubeBootstrapperNamespaceFromEnv
		}

		if !alphaNumDash.MatchString(kubeBootstrapperNamespace) {
			return fmt.Errorf("%s is not a valid namespace", kubeBootstrapperNamespace)
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&kubeBootstrapperNamespace, "namespace", "n", defaultKubeBootstrapperNamespace, "Namespace in which kubebootstrapper is installed [$KUBEBOOTSTRAPPER_NAMESPACE]")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Turn on debug logging")
	RootCmd.AddCommand(newCmdInstall())
	RootCmd.AddCommand(newCmdUninstall())
	RootCmd.AddCommand(newCmdVersion())
}
