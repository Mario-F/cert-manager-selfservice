package kube

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type CertifcateResult struct {
	Domain     string
	CertsFound []KubeCertificate
}

type KubeCertificate struct {
	Certificate certv1.Certificate
	Domains     []string
	Secret      corev1.Secret
	LastAccess  int64
	Ready       bool
}

var (
	managerId          string
	promCertGauge      prometheus.Gauge
	promCertReadyGauge *prometheus.GaugeVec
)

func init() {
	promCertGauge = promauto.NewGauge(prometheus.GaugeOpts{Name: "cms_certificates", Help: "The amount of managed certificates"})
	promCertReadyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "cms_certificate_ready", Help: "Shows ready status of a certificate"}, []string{"domain"})
}

func SetManagerId(newId string) {
	managerId = newId
}

func GetCertificates() ([]KubeCertificate, error) {
	log.Debugf("Get all cert-manager-selfservice managed certificates with ID: %s", managerId)
	result := []KubeCertificate{}

	client, err := getClient("")
	if err != nil {
		return result, err
	}

	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"cert-manager-selfservice/managed": managerId}}
	listOptions := metav1.ListOptions{LabelSelector: labels.Set(labelSelector.MatchLabels).String()}
	kubeResult, err := client.CertManager.CertmanagerV1().Certificates(client.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return result, err
	}

	promCertGauge.Set(float64(len(kubeResult.Items)))
	log.Debugf("Found %d certificates for manager-id %s", len(kubeResult.Items), managerId)
	for _, c := range kubeResult.Items {
		actCert := KubeCertificate{Ready: true}

		actCert.Domains = append(actCert.Domains, c.Spec.DNSNames...)

		// now get tls secret
		secret, err := client.K8s.CoreV1().Secrets(client.Namespace).Get(context.TODO(), c.Spec.SecretName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("TLS Secret for domain %s not ready yet: %s", c.Name, c.Spec.SecretName)
			actCert.Ready = false
		}
		actCert.Certificate = c
		actCert.Secret = *secret

		// TODO: Add prom metrics (last access time gauge)
		actCert.LastAccess = parseTime(c.Name, c.ObjectMeta.Labels["cert-manager-selfservice/last-access"])

		if actCert.Ready {
			promCertReadyGauge.WithLabelValues(actCert.Domains[0]).Set(1)
		} else {
			promCertReadyGauge.WithLabelValues(actCert.Domains[0]).Set(0)
		}

		result = append(result, actCert)
	}

	return result, nil
}

func (k KubeCertificate) updateAccess() error {
	log.Debugf("Update lastAccess time for %s", k.Certificate.Name)

	client, err := getClient("")
	if err != nil {
		return err
	}

	timeNow := time.Now().Unix()
	k.Certificate.ObjectMeta.Labels["cert-manager-selfservice/last-access"] = fmt.Sprintf("%d", timeNow)
	_, err = client.CertManager.CertmanagerV1().Certificates(k.Certificate.Namespace).Update(context.TODO(), &k.Certificate, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k KubeCertificate) Delete() error {
	log.Debugf("Delete certificate %s", k.Certificate.Name)

	client, err := getClient("")
	if err != nil {
		return err
	}

	err = client.CertManager.CertmanagerV1().Certificates(k.Certificate.Namespace).Delete(context.TODO(), k.Certificate.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetCertificate(domain string, updateAccess bool) (CertifcateResult, error) {
	log.Infof("Search for domain %s certificate", domain)
	result := CertifcateResult{Domain: domain}

	cCerts, err := GetCertificates()
	if err != nil {
		return result, err
	}

	// Search for domain in dns names
	var kCerts []KubeCertificate
	for _, kc := range cCerts {
		for _, d := range kc.Domains {
			if d == domain {
				kCerts = append(kCerts, kc)

				// TODO: Add prom metrics (access count)
				if updateAccess {
					err := kc.updateAccess()
					if err != nil {
						return result, err
					}
				}
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

func CreateCertificate(domain string, issuer cmmeta.ObjectReference) error {
	log.Infof("Create certificate for domain %s", domain)
	domainSlug := strings.ReplaceAll(domain, ".", "-")

	client, err := getClient("")
	if err != nil {
		return err
	}

	timeNow := time.Now().Unix()
	crt := &certv1.Certificate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Certificate",
			APIVersion: "cert-manager.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", managerId, domainSlug),
			Labels: map[string]string{
				"cert-manager-selfservice/managed":     managerId,
				"cert-manager-selfservice/last-access": fmt.Sprintf("%d", timeNow),
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
	// TODO: Add prom metrics (create count)

	return nil
}
