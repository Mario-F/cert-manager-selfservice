package kube

import (
	"context"
	"fmt"
	"strings"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
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

func CreateCertificate(domain string, issuer cmmeta.ObjectReference) error {
	log.Infof("Create certificate for domain %s", domain)
	domainSlug := strings.ReplaceAll(domain, ".", "-")

	client, err := getClient("")
	if err != nil {
		return err
	}

	crt := &certv1.Certificate{
		TypeMeta: v1.TypeMeta{
			Kind:       "Certificate",
			APIVersion: "cert-manager.io/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("cert-%s", domainSlug),
			Labels: map[string]string{
				"cert-manager-selfservice/managed": "true",
			},
		},
		Spec: certv1.CertificateSpec{
			CommonName: domain,
			DNSNames: []string{
				domain,
			},
			SecretName: fmt.Sprintf("tls-%s", domainSlug),
			IssuerRef:  issuer,
			Usages: []certv1.KeyUsage{
				"server auth",
				"client auth",
			},
		},
	}

	c, err := client.CertManager.CertmanagerV1().Certificates(client.Namespace).Create(context.TODO(), crt, v1.CreateOptions{})
	if err != nil {
		return err
	}
	_ = c

	return nil
}
