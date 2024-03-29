package kube

import (
	"os"
	"path/filepath"
	"sync"

	discovery "github.com/gkarthiks/k8s-discovery"
	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeClients struct {
	Namespace   string
	Version     string
	K8s         *kubernetes.Clientset
	CertManager *cmclient.Clientset
	Created     bool
}

var (
	kubeConfigPath    string
	fetchedClients    KubeClients
	fetchLock         sync.Mutex
	kubeNamespace     string
	promKubeApiAccess prometheus.Counter
)

func init() {
	promKubeApiAccess = promauto.NewCounter(prometheus.CounterOpts{Name: "cms_kube_api_access_total", Help: "Count the number of kubernetes api calls"})
}

func SetKubeConfigPath(path string) {
	kubeConfigPath = path
}

func SetNamespace(kubens string) {
	log.Infof("Setting kubernetes namespace to %s", kubens)
	kubeNamespace = kubens
}

func getClient() (KubeClients, error) {
	log.Debug("Get kube client by trying ClusterConfig")
	promKubeApiAccess.Inc()

	fetchLock.Lock()
	defer fetchLock.Unlock()

	if fetchedClients.Created {
		log.Debug("Kube clients already exists")
		return fetchedClients, nil
	}

	result := KubeClients{}

	// Try get config from passed filepath
	if kubeConfigPath == "" {
		// check if env variable KUBECONFIG is set
		if kubeConfigPath = os.Getenv("KUBECONFIG"); kubeConfigPath == "" {
			// if not set, try to get config from home directory
			kubeConfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}
	}
	log.Debugf("Kubeconfig path is set to: %s", kubeConfigPath)
	// Set KUBECONFIG path env variable as workaround for k8s-discovery
	os.Setenv("KUBECONFIG", kubeConfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Debugf("ClusterConfig created error %v+", err.Error())
		// Try inclusterconfig
		config, err = rest.InClusterConfig()
		if err != nil {
			return result, err
		}
	}

	// Get k8s and cert-manager clients
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

	// Add additional information
	k8s, err := discovery.NewK8s()
	if err != nil {
		return result, err
	}
	namespace, err := k8s.GetNamespace()
	if err != nil {
		return result, err
	}
	if namespace == "" {
		result.Namespace = kubeNamespace
	} else {
		result.Namespace = namespace
	}
	version, err := k8s.GetVersion()
	if err != nil {
		return result, err
	}
	result.Version = version

	result.Created = true
	fetchedClients = result

	return result, nil
}

func CheckKubernetes() error {
	_, err := getClient()
	return err
}
