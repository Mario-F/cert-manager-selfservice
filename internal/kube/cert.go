package kube

import (
	"context"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertifcateResult struct {
	Domain     string
	CertsFound []certv1.Certificate
}

func GetCertificate(domain string) (CertifcateResult, error) {
	log.Infof("Search for domain %s certificate", domain)
	result := CertifcateResult{Domain: domain}

	client, err := getClient("")
	if err != nil {
		return result, err
	}

	cCerts, err := client.CertManager.CertmanagerV1().Certificates(client.Namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return result, err
	}

	// Search for domain in dns names
	var fCerts []certv1.Certificate
	for _, c := range cCerts.Items {
		for _, d := range c.Spec.DNSNames {
			if d == domain {
				fCerts = append(fCerts, c)
				break
			}
		}
	}
	result.CertsFound = fCerts

	return result, nil
}
