package kube

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertifcateResult struct {
	Domain     string
	CertsFound []KubeCertificate
}

type KubeCertificate struct {
	Certificate certv1.Certificate
	Secret      corev1.Secret
	LastAccess  int64
	Ready       bool
}

func GetCertificate(domain string, updateAccess bool) (CertifcateResult, error) {
	log.Infof("Search for domain %s certificate", domain)
	result := CertifcateResult{Domain: domain}

	client, err := getClient("")
	if err != nil {
		return result, err
	}

	cCerts, err := client.CertManager.CertmanagerV1().Certificates(client.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	// Search for domain in dns names
	var kCerts []KubeCertificate
	for _, c := range cCerts.Items {
		for _, d := range c.Spec.DNSNames {
			if d == domain {
				// Certificate found, now get tls secret
				kCert := KubeCertificate{Ready: true}
				secret, err := client.K8s.CoreV1().Secrets(client.Namespace).Get(context.TODO(), c.Spec.SecretName, metav1.GetOptions{})
				if err != nil {
					log.Errorf("TLS Secret for domain %s not ready yet: %s", domain, c.Spec.SecretName)
					kCert.Ready = false
				}
				kCert.Certificate = c
				kCert.Secret = *secret
				kCerts = append(kCerts, kCert)

				// Update timestamp on last access label
				if updateAccess {
					timeNow := time.Now().Unix()
					timeCert := parseTime(domain, c.ObjectMeta.Labels["cert-manager-selfservice/last-access"])

					if timeCert < (timeNow - 5) {
						log.Debugf("Update lastAccess time for domain %s", domain)
						c.ObjectMeta.Labels["cert-manager-selfservice/last-access"] = fmt.Sprintf("%d", timeNow)
						_, err := client.CertManager.CertmanagerV1().Certificates(c.Namespace).Update(context.TODO(), &c, metav1.UpdateOptions{})
						if err != nil {
							return result, err
						}
					}
				}

				// Added lastAccess
				kCert.LastAccess = parseTime(domain, c.ObjectMeta.Labels["cert-manager-selfservice/last-access"])

				break
			}
		}
	}
	result.CertsFound = kCerts

	return result, nil
}

func parseTime(domain string, stringTime string) int64 {
	iTimestamp, err := strconv.ParseInt(stringTime, 10, 64)
	log.Debugf("Last access timestamp for domain %s is %d", domain, iTimestamp)
	if err != nil {
		log.Warnf("Failed to parse lastAccess for cert with domain %s", domain)
		return 0
	}
	return iTimestamp
}

func CreateCertificate(domain string, issuer cmmeta.ObjectReference, certPrefix string) error {
	log.Infof("Create certificate for domain %s", domain)
	domainSlug := strings.ReplaceAll(domain, ".", "-")

	client, err := getClient("")
	if err != nil {
		return err
	}

	crt := &certv1.Certificate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Certificate",
			APIVersion: "cert-manager.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", certPrefix, domainSlug),
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

	c, err := client.CertManager.CertmanagerV1().Certificates(client.Namespace).Create(context.TODO(), crt, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	_ = c

	return nil
}
