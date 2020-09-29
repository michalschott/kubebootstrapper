# kubebootstrap

Aim of this project is to mainly upskill in golang (k8s sdk, API etc).

There are varous flavours of k8s available, and wheel is reinvented when new project kicks in. What to use to bootstrap my cluster - helm? kustomize? tanka?

This project aims to make bootstrap process easier (k8s flavour doesn't matter) - with a controller which will maintain cluster bootstrapped state as desired/specified in [configmap/config file]. Also updates should seamless (in love by CoreOS/Flatcar linux update model, would like to see if this would work for k8s).

Project is developed during free time (obviously this is not infinite resource), will see where it goes...

# cli

CLI to install/uninstall controller into your k8s cluster. Inspired by [linkerd2-cli](https://github.com/linkerd/linkerd2/tree/main/cli).

`CLI install` should be run by `cluster-admin user`.

# controller

Simple control-loop style application which will fetch manifests from [S3-compatible place, alpha/beta/stable channel], and based on the user-configuration apply part of them.

# dev

Requirements: `kind create cluster`
