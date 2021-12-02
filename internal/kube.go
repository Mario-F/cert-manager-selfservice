package kube

import (
	"path/filepath"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeClients struct {
	K8s         *kubernetes.Clientset
	CertManager *cmclient.Clientset
}

func getClient(kubeConfigPath string) (KubeClients, error) {
	log.Debug("Get kube client by trying ClusterConfig")
	result := KubeClients{}

	// Try get config from passed filepath
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}
	log.Debugf("Kubeconfig path is set to: %s", kubeConfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Debugf("ClusterConfig created error %v+", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return result, err
	}
	result.K8s = clientset

	certClientset, err := cmclient.NewForConfig(config)
	if err != nil {
		return result, err
	}
	result.CertManager = certClientset

	return result, nil
}
