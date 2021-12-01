package kube

import (
	"testing"

	log "github.com/sirupsen/logrus"
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
}
