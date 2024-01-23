package web

import (
	"sync"
	"time"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509/pkix"
	"crypto/x509"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

// selfSignedCertificateGetter returns a function that can be used as a certificate getter in TLSConfig.
func selfSignedCertificateGetter() (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error) {
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating key: %v", err)
	}

	rootKeyPEM, err := encodeKeyAsPEM(rootKey)
	if err != nil {
		return nil, err
	}

	serialNumber := big.NewInt(0)

	tmpl := &x509.Certificate{
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Neds International"}},
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		BasicConstraintsValid: true,
		IPAddresses:           nil,
	}

	generateCertificate := func() (*tls.Certificate, error) {
		tmpl.SerialNumber.Add(tmpl.SerialNumber, big.NewInt(1))

		now := time.Now()
		tmpl.NotBefore = now
		tmpl.NotAfter = now.Add(1 * time.Hour)

		certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, rootKey.Public(), rootKey)
		if err != nil {
			return nil, err
		}

		b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
		rootCertPEM := pem.EncodeToMemory(&b)

		certificate, err := tls.X509KeyPair(rootCertPEM, rootKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("invalid key pair: %v", err)
		}

		return &certificate, nil
	}

	certificate, err := generateCertificate()

	var m sync.Mutex

	return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		lockTime := time.Now()

		m.Lock()

		if time.Now().After(tmpl.NotAfter) {
			var log logrus.FieldLogger = logrus.StandardLogger()

			defer func(genTime time.Time) {
				log.WithFields(logrus.Fields{
					"dur_seconds":           time.Since(genTime).Seconds(),
					"dur_with_lock_seconds": time.Since(lockTime).Seconds(),
					"not_before":            tmpl.NotBefore.Format(time.RFC3339),
					"not_after":             tmpl.NotAfter.Format(time.RFC3339),
				}).Debug("Service certificate regeneration complete")
			}(time.Now())

			certificate, err = generateCertificate()
		}

		m.Unlock()

		return certificate, err
	}, err
}

// encodeKeyAsPEM encodes the given private key in PEM format.
func encodeKeyAsPEM(key crypto.PrivateKey) (rootKeyPEM []byte, err error) {
	var block pem.Block

	switch k := key.(type) {
	case *rsa.PrivateKey:
		block.Type = "RSA PRIVATE KEY"
		block.Bytes = x509.MarshalPKCS1PrivateKey(k)
	case *ecdsa.PrivateKey:
		block.Type = "EC PRIVATE KEY"
		block.Bytes, err = x509.MarshalECPrivateKey(k)
	}

	if err != nil {
		return nil, err
	}

	rootKeyPEM = pem.EncodeToMemory(&block)

	return rootKeyPEM, err
}

