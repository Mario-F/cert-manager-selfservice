package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/logger"
)

type CaCert struct {
	Cert    *x509.Certificate
	PEM     bytes.Buffer
	PrivPEM bytes.Buffer
	PrivKey *rsa.PrivateKey
	Created bool
}

type SelfSignedCertRequest struct {
	Domain string
}

type SelfSignedCertResult struct {
	ServerTLSConf *tls.Config
	ServerCert    tls.Certificate
	ClientTLSConf *tls.Config
	CaCertPEM     bytes.Buffer
	CertPEM       bytes.Buffer
	PrivatePEM    bytes.Buffer
}

var ca CaCert

func SelfSignedCert(req SelfSignedCertRequest) (SelfSignedCertResult, error) {
	ca, err := getCA()
	logger.Verbosef("Certficate with domain %s was requested", req.Domain)
	result := SelfSignedCertResult{}
	if err != nil {
		return result, err
	}
	result.CaCertPEM = ca.PEM
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Selfsigned"},
			Country:       []string{"XX"},
			Province:      []string{""},
			Locality:      []string{"Global"},
			StreetAddress: []string{""},
			PostalCode:    []string{"99999"},
			CommonName:    req.Domain,
		},
		DNSNames:     []string{req.Domain},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return result, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &certPrivKey.PublicKey, ca.PrivKey)
	if err != nil {
		return result, err
	}
	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return result, err
	}
	result.CertPEM = *certPEM
	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return result, err
	}
	result.PrivatePEM = *certPrivKeyPEM
	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return result, err
	}
	result.ServerCert = serverCert

	result.ServerTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca.PEM.Bytes())
	result.ClientTLSConf = &tls.Config{
		RootCAs: certpool,
	}

	logger.Verbosef("Certificate for domain %s finished", req.Domain)
	return result, nil
}

// getCa return a fresh or already created self signed ca
func getCA() (*CaCert, error) {
	if ca.Created {
		logger.Debugf("CA is already created, return existing")
		return &ca, nil
	}
	logger.Debugf("CA not created, generating CA-Cert")
	caCertData := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Selfsigned"},
			Country:       []string{"XX"},
			Province:      []string{""},
			Locality:      []string{"Global"},
			StreetAddress: []string{""},
			PostalCode:    []string{"99999"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	ca.Cert = caCertData
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return &ca, err
	}
	ca.PrivKey = caPrivKey
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertData, caCertData, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return &ca, err
	}
	err = pem.Encode(&ca.PEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return &ca, err
	}

	err = pem.Encode(&ca.PrivPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return &ca, err
	}

	ca.Created = true
	return &ca, nil
}
