package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func newCmdUninstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Args:  cobra.NoArgs,
		Short: "Uninstalls/updates kubebootstrapper controller",
		Long:  "Uninstalls/updates kubebootstrapper controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstallRunE(cmd.Context(), stdout)
		},
	}
	return cmd
}

func uninstallRunE(ctx context.Context, wout io.Writer) error {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Writer = wout
	spinner.Prefix = "Installing: "
	spinner.Color("bold")

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		spinner.Stop()
		fmt.Println("Uninstallation NOT OK " + failStatus + "\n")
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		spinner.Stop()
		fmt.Println("Uninstallation NOT OK " + failStatus + "\n")
		return err
	}

	err = deleteDeployment(ctx, clientset)
	if err != nil {
		spinner.Stop()
		fmt.Println("Uninstallation NOT OK " + failStatus + "\n")
		return err
	}

	err = deleteNamespace(ctx, clientset)
	if err != nil {
		spinner.Stop()
		fmt.Println("Uninstallation NOT OK " + failStatus + "\n")
		return err
	}

	spinner.Stop()
	fmt.Println("Uninstallation OK " + okStatus)
	return nil
}

func deleteDeployment(ctx context.Context, clientset *kubernetes.Clientset) error {
	client := clientset.AppsV1().Deployments(kubeBootstrapperNamespace)
	err := client.Delete(ctx, "controller", metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

func deleteNamespace(ctx context.Context, clientset *kubernetes.Clientset) error {
	client := clientset.CoreV1().Namespaces()
	err := client.Delete(ctx, kubeBootstrapperNamespace, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}
