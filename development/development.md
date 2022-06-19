# Development

This folder will contain all the files that are needed for advanced development and testing but not mendatory to use.

## Cluster

The `./cluster.sh` script can used to spin up a k3s cluster with cert-manager and a self signing clusterissuer.

Requirements on host:

* shell / linux (wsl-ubuntu tested but mac os can also work)
* docker
* kubectl
* helm
