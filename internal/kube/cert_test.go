package kube

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestKubeCertHandling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	fDomain := "example.com"
	t.Run("Test get domain", func(t *testing.T) {
		cert, err := GetCertificate(fDomain)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Got cert %v+", cert)
	})
}
