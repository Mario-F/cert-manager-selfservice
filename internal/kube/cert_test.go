package kube

import (
	"testing"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	log "github.com/sirupsen/logrus"
)

func TestKubeCertHandling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	fDomain := "testcert.1234.com"
	t.Run("Test get domain", func(t *testing.T) {
		cert, err := GetCertificate(fDomain, true)
		if err != nil {
			t.Logf("Test can fail %v+", err)
			return
		}
		t.Logf("Got cert %v+", cert)
	})

	t.Run("Test create domain", func(t *testing.T) {
		testIssuer := cmmeta.ObjectReference{
			Name: "TestClusterIssuer",
			Kind: "ClusterIssuer",
		}
		err := CreateCertificate(fDomain, testIssuer, "test")
		if err != nil {
			t.Logf("Test can fail %v+", err)
			return
		}
	})
}
