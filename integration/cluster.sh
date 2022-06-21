#!/bin/sh
# Spin up a kubernetes cluster with docker and install cert-manager with self signing certs

K3S_VERSION=${K3S_VERSION:-v1.22.10-k3s1}
CERT_MANAGER_VERSION=${CERT_MANAGER_VERSION:-v1.8.1}
HOST_KUBECONFIG=${HOST_KUBECONFIG:-/tmp/cms-kubeconfig}
BASEDIR=$(dirname "$0")
DEBUG=${DEBUG:-false}

ARG_COMMAND=${1:-}

# check requirements
command -v docker >/dev/null 2>&1 || { echo >&2 "docker is required but not installed.  Aborting."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is required but not installed.  Aborting."; exit 1; }
command -v helm >/dev/null 2>&1 || { echo >&2 "helm is required but not installed.  Aborting."; exit 1; }

# switch case for the command
case $ARG_COMMAND in
  "start")
    # run k3s in docker
    echo "start k3s in docker..."
    docker run -d --privileged --name=cms-k3s \
      -e K3S_KUBECONFIG_OUTPUT=/tmp/kubeconfig \
      -p 6443:6443 \
      rancher/k3s:$K3S_VERSION server
    if [ $? -ne 0 ]; then
      echo "failed to start k3s in docker, please check error message"
      exit 1
    fi
    if [ $DEBUG = true ]; then
      docker ps
    fi

    # wait for k3s to be ready
    echo "wait for k3s to be ready..."
    K3S_READY=1
    while [ $K3S_READY -ne 0 ]; do
      sleep 2
      docker exec -it cms-k3s kubectl get nodes | grep Ready
      K3S_READY=$?
      if [ $DEBUG = true ]; then
        echo "exit code: $K3S_READY"
        docker exec -i cms-k3s kubectl get nodes
      fi
    done
    echo "k3s is ready"
    if [ $DEBUG = true ]; then
      docker ps
    fi

    # copy kubeconfig to host
    echo "copy kubeconfig to host..."
    docker cp cms-k3s:/tmp/kubeconfig $HOST_KUBECONFIG

    # install cert-manager
    echo "install cert-manager..."
    helm install \
      --wait \
      --kubeconfig $HOST_KUBECONFIG \
      --repo https://charts.jetstack.io \
      --namespace cert-manager \
      --create-namespace \
      --version $CERT_MANAGER_VERSION \
      --values ${BASEDIR}/cluster/cert-manager.values.yaml \
      cert-manager cert-manager
    if [ $? -ne 0 ]; then
      echo "failed to install cert-manager, please check error message"
      exit 1
    fi
    echo "cert-manager is installed"
    sleep 2 # webhook possible not started yet
    echo "install cluster-issuer..."
    kubectl --kubeconfig=$HOST_KUBECONFIG \
      --namespace cert-manager \
      apply -f ${BASEDIR}/cluster/cluster-issuer.yaml
    if [ $? -ne 0 ]; then
      echo "failed to install cluster issuer, please check error message"
      exit 1
    fi
    echo "cluster-issuer is installed"

    kubectl --kubeconfig=$HOST_KUBECONFIG create namespace cms

    echo ""
    echo "use kubeconfig from host with:"
    echo "export KUBECONFIG=${HOST_KUBECONFIG}"
    ;;
  "stop")
    # stop and remove container
    echo "stop and remove container..."
    docker stop cms-k3s
    docker rm cms-k3s
    echo "container stopped and removed"
    ;;
  "usage")
    echo "usage: $0 [start|stop]"
    ;;
esac
