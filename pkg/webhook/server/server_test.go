package server

import (
	"crypto/x509"
	"github.com/go-chi/chi"
	"io/ioutil"
	"testing"
)

func Test_ServerStart(t *testing.T) {
	s := &Settings{
		ServerCertPath: "testdata/demo-certs/server.crt",
		ServerKeyPath:  "testdata/demo-certs/server-key.pem",
	}

	rtr := chi.NewRouter()

	srv := &WebhookServer{
		Settings: s,
		Router:   rtr,
	}

	err := srv.Start()
	if err != nil {
		t.Fatalf("Server should start: %v", err)
	}
}

func Test_Client_CA(t *testing.T) {
	roots := x509.NewCertPool()

	s := Settings{}
	s.ClientCAPaths = []string{
		"testdata/demo-certs/client-ca.pem",
	}

	for _, caPath := range s.ClientCAPaths {
		caBytes, err := ioutil.ReadFile(caPath)
		if err != nil {
			t.Fatalf("ca '%s' should be read: %v", caPath, err)
		}

		ok := roots.AppendCertsFromPEM(caBytes)
		if !ok {
			t.Fatalf("ca '%s' should be parsed", caPath)
		}
	}
}
