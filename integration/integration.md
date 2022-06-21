# Integration

This folder will contain all the files that are needed for an full integration test.

## Cluster

The `./cluster.sh` script can used to spin up a k3s cluster with cert-manager and a self signing clusterissuer.

Requirements on host:

* shell / linux (wsl-ubuntu tested but mac os can also work)
* docker
* kubectl
* helm

Start with `./integration/cluster.sh start`

Develop with `./debug server --kubeconfig /tmp/cms-kubeconfig --kube-namespace cms --issuer-name cms-development-cluster-issuer`

Stop with `./integration/cluster.sh stop`
