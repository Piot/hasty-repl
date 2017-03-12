package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func hackServerCertAsCA(cacert string) (*x509.CertPool, error) {
	caPool := x509.NewCertPool()
	severCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load CA cert: '%s'", err)
	}
	caPool.AppendCertsFromPEM(severCert)
	return caPool, nil
}

// NewConnection : Creates a new connection
func NewConnection(host string, cacert string) (*tls.Conn, error) {
	var caPool *x509.CertPool

	var tlsConfig *tls.Config

	if cacert != "" {
		var caErr error
		caPool, caErr = hackServerCertAsCA(cacert)
		if caErr != nil {
			return nil, caErr
		}
		tlsConfig = &tls.Config{KeyLogWriter: os.Stdout, RootCAs: caPool, MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS12}
	} else {
		caPool, caErr := x509.SystemCertPool()
		if caErr != nil {
			return nil, caErr
		}
		tlsConfig = &tls.Config{ClientCAs: caPool, RootCAs: caPool,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
	}
	conn, err := tls.Dial("tcp", host, tlsConfig)
	if err != nil {

		return nil, fmt.Errorf("Could not connect '%s'", err)
	}
	log.Printf("Connected to '%s'", host)
	// defer conn.Close()

	return conn, nil
}
