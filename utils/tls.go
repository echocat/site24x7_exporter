package utils

import (
	"crypto/x509"
	"io/ioutil"
)

func LoadInternalCaBundle() *x509.CertPool {
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM([]byte(caBundle))
	return certificates
}

func LoadCertificatesFrom(pemFile string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM(caCert)
	return certificates, nil
}


