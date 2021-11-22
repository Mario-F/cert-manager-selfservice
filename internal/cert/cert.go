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
	cert          *x509.Certificate
	caCertPEM     bytes.Buffer
	caCertPrivPEM bytes.Buffer
	caPrivKey     *rsa.PrivateKey
	created       bool
}

type SelfSignedCertRequest struct {
	domain string
}

type SelfSignedCertResult struct {
	ServerTLSConf *tls.Config
	ClientTLSConf *tls.Config
}

var ca CaCert

func SelfSignedCert(req SelfSignedCertRequest) (SelfSignedCertResult, error) {
	ca, err := getCA()
	result := SelfSignedCertResult{}
	if err != nil {
		return result, err
	}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Selfsigned"},
			Country:       []string{"XX"},
			Province:      []string{""},
			Locality:      []string{"Global"},
			StreetAddress: []string{""},
			PostalCode:    []string{"99999"},
		},
		DNSNames:     []string{req.domain},
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
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.cert, &certPrivKey.PublicKey, ca.caPrivKey)
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
	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return result, err
	}
	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return result, err
	}

	result.ServerTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca.caCertPEM.Bytes())
	result.ClientTLSConf = &tls.Config{
		RootCAs: certpool,
	}

	return result, nil
}

// getCa return a fresh or already created self signed ca
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
	ca.cert = caCertData
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return &ca, err
	}
	ca.caPrivKey = caPrivKey
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