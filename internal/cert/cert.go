package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/logger"
)

type CaCert struct {
	caCertPEM     bytes.Buffer
	caCertPrivPEM bytes.Buffer
	created       bool
}

var ca CaCert

func SelfSignedCert() error {
	_, err := getCA()
	if err != nil {
		return err
	}
	return nil
}

func getCA() (*CaCert, error) {
	if ca.created {
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
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return &ca, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertData, caCertData, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return &ca, err
	}
	err = pem.Encode(&ca.caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return &ca, err
	}

	err = pem.Encode(&ca.caCertPrivPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return &ca, err
	}

	ca.created = true
	return &ca, nil
}
