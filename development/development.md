# Development

This folder will contain all the files that are needed for advanced development and testing but not mendatory to use.

## Cluster

The `./cluster.sh` script can used to spin up a k3s cluster with cert-manager and a self signing clusterissuer.

Requirements on host:

* shell / linux (wsl-ubuntu tested but mac os can also work)
* docker
* kubectl
* helm

Start with `./development/cluster.sh start`

Develop with `./debug server --kubeconfig /tmp/cms-kubeconfig --kube-namespace cms --issuer-name cms-development-cluster-issuer`

Stop with `./development/cluster.sh stop`
