package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math"
	"math/big"
	"time"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
)

func GenerateKeystore(
	storepass,
	keypass,
	alias,
	cname string,

	validity time.Duration,
	bits uint32,

	writer io.Writer,
) error {
	// Generate private key
	key, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return err
	}

	rawKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	// Generate certificate
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return err
	}

	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		NotBefore:    now,
		NotAfter:     now.Add(validity),
		Subject: pkix.Name{
			CommonName: cname,
		},
		Issuer: pkix.Name{
			CommonName: cname,
		},
	}

	cert, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		return err
	}

	// Generate & write keystore
	ks := keystore.New()

	if err := ks.SetPrivateKeyEntry(
		alias,
		keystore.PrivateKeyEntry{
			CreationTime: time.Now(),
			PrivateKey:   rawKey,
			CertificateChain: []keystore.Certificate{
				{
					Type:    "X509",
					Content: cert,
				},
			},
		},
		[]byte(keypass),
	); err != nil {
		return err
	}

	return ks.Store(writer, []byte(storepass))
}
