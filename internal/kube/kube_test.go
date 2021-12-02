package kube

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubeClientConfig(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("Test creating kube client", func(t *testing.T) {
		_, err := getClient("")
		if err != nil {
			t.Log("Kube client error but ok in pipeline test")
		} else {
			t.Log("Kube client success maybe on local development")
		}
	})

	t.Run("Test get certificates", func(t *testing.T) {
		client, err := getClient("")
		if err != nil {
			t.Log("Kube client error but ok in pipeline test")
			return
		}
		pods, _ := client.K8s.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
		t.Logf("%v+", pods)
		certs, _ := client.CertManager.CertmanagerV1().Certificates("unraid").List(context.TODO(), v1.ListOptions{})
		t.Logf("%v+", certs)
	})
}
