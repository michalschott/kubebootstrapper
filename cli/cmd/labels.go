package cmd

import "github.com/michalschott/kubebootstrapper/pkg/version"

func generateLabels() map[string]string {
	labels := map[string]string{
		"version": version.Version,
	}

	return labels
}
