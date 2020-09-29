package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/briandowns/spinner"
	"github.com/michalschott/kubebootstrapper/pkg/ptr"
	version "github.com/michalschott/kubebootstrapper/pkg/version"
)

type installOptions struct {
	kubeBootstrapperVersion string
	kubernetesVersion       string
	kubernetesFlavour       string
}

// newInstallOptionsWithDefaults initializes install options with default options
func newInstallOptionsWithDefaults() *installOptions {
	return &installOptions{
		kubeBootstrapperVersion: version.Version,
		kubernetesVersion:       "1.17.10",
		kubernetesFlavour:       "EKS",
	}
}

func newCmdInstall() *cobra.Command {
	options := newInstallOptionsWithDefaults()

	cmd := &cobra.Command{
		Use:   "install",
		Args:  cobra.NoArgs,
		Short: "Installs/updates kubebootstrapper controller",
		Long:  "Installs/updates kubebootstrapper controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installRunE(cmd.Context(), options, stdout)
		},
	}
	return cmd
}

func installRunE(ctx context.Context, options *installOptions, wout io.Writer) error {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Writer = wout
	spinner.Prefix = "Installing: "
	spinner.Color("bold")

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		spinner.Stop()
		fmt.Println("Installation NOT OK " + failStatus + "\n")
		return err
	}

	err = createNamespace(ctx, clientset)
	if err != nil {
		spinner.Stop()
		fmt.Println("Installation NOT OK " + failStatus + "\n")
		return err
	}

	err = createDeployment(ctx, clientset)
	if err != nil {
		spinner.Stop()
		fmt.Println("Installation NOT OK " + failStatus + "\n")
		return err
	}

	spinner.Stop()
	fmt.Println("Installation OK " + okStatus)
	return nil
}

func createNamespace(ctx context.Context, clientset *kubernetes.Clientset) error {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   kubeBootstrapperNamespace,
			Labels: generateLabels(),
		},
	}

	client := clientset.CoreV1().Namespaces()
	_, err := client.Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			_, err = client.Update(ctx, namespace, metav1.UpdateOptions{})
		} else {
			return err
		}
	}

	return nil
}

func createDeployment(ctx context.Context, clientset *kubernetes.Clientset) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller",
			Namespace: kubeBootstrapperNamespace,
			Labels:    generateLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: generateLabels(),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: generateLabels(),
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "controller",
							Image: "controller:" + version.Version,
						},
					},
				},
			},
		},
	}

	if verbose {
		deployment.Spec.Template.Spec.Containers[0].Args = []string{"-l", "debug"}
	}

	client := clientset.AppsV1().Deployments(deployment.Namespace)
	_, err := client.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			_, err = client.Update(ctx, deployment, metav1.UpdateOptions{})
		} else {
			return err
		}
	}

	return nil
}
