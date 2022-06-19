#!/bin/sh
# Spin up a kubernetes cluster with docker and install cert-manager with self signing certs

K3S_VERSION=${K3S_VERSION:-v1.22.10-k3s1}
HOST_KUBECONFIG=${HOST_KUBECONFIG:-/tmp/cms-kubeconfig}

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

    # wait for k3s to be ready
    echo "wait for k3s to be ready..."
    K3S_READY=1
    while [ $K3S_READY -ne 0 ]; do
      sleep 1
      docker exec -it cms-k3s kubectl get nodes &> /dev/null
      K3S_READY=$?
    done
    echo "k3s is ready"

    # copy kubeconfig to host
    echo "copy kubeconfig to host..."
    docker cp cms-k3s:/tmp/kubeconfig $HOST_KUBECONFIG

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
