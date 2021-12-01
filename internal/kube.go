package kube

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getClient(kubeConfigPath string) (*kubernetes.Clientset, error) {
	log.Debug("Get kube client by trying ClusterConfig")

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
		return clientset, err
	}

	return clientset, nil
}
