package rsahelper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/apex/log"
)

func RsaSetup(rsaPrivateKeyLocation string) (string, error) {

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		log.Fatalf("RSA private key file not found")
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte

	if privPem.Type != "RSA PRIVATE KEY" {
		log.Fatalf("RSA private key is of the wrong type")
	}

	privPemBytes = privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		log.Fatalf("Unable to parse RSA private key")
	}

	var ok bool
	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatalf("Unable to parse RSA private key")
	}

	privatKeyString := ExportRsaPrivateKeyAsPemStr(privateKey)

	return privatKeyString, nil
}

func GenRSA(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalf("Failed to generate signing key", err)
	}
	return key, err
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func GetRSA(bits int) (string, error) {
	key, err := GenRSA(bits)
	if err != nil {
		log.Errorf("RSA Generation failed")
	}
	return ExportRsaPrivateKeyAsPemStr(key), err
}
