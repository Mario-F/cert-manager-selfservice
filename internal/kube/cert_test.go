package kube

import (
	"testing"

	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	log "github.com/sirupsen/logrus"
)

func TestDomainToSlug(t *testing.T) {
	tests := []struct {
		domain   string
		expected string
	}{
		{"example.com", "example-com"},
		{"sub.example.com", "sub-example-com"},
		{"*.example.com", "wildcard-example-com"},
		{"*.sub.example.com", "wildcard-sub-example-com"},
	}
	for _, tt := range tests {
		result := domainToSlug(tt.domain)
		if result != tt.expected {
			t.Errorf("domainToSlug(%q) = %q, want %q", tt.domain, result, tt.expected)
		}
	}
}

func TestKubeCertHandling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	fDomain := "testcert.1234.com"
	t.Run("Test get domain", func(t *testing.T) {
		cert, err := GetCertificate(fDomain, true, false)
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
		err := CreateCertificate(fDomain, testIssuer)
		if err != nil {
			t.Logf("Test can fail %v+", err)
			return
		}
	})
}
