package kube

import "k8s.io/client-go/rest"

func Client() {
	_, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
}
