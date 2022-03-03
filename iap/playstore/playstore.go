package playstore

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
)

// New returns http client which includes the credentials to access androidpublisher API.
func New(jsonKey []byte) (*Client, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	val := conf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{service}, err
}

// NewWithClient returns http client which includes the custom http client.
func NewWithClient(jsonKey []byte, cli *http.Client) (*Client, error) {
	if cli == nil {
		return nil, fmt.Errorf("client is nil")
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cli)

	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{service}, err
}

// VerifySignature verifies in app billing signature.
// You need to prepare a public key for your Android app's in app billing
// at https://play.google.com/apps/publish/
func VerifySignature(base64EncodedPublicKey string, receipt []byte, signature string) (isValid bool, err error) {
	// prepare public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key")
	}
	publicKey, _ := publicKeyInterface.(*rsa.PublicKey)

	// generate hash value from receipt
	hasher := sha1.New()
	hasher.Write(receipt)
	hashedReceipt := hasher.Sum(nil)

	// decode signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature")
	}

	// verify
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature); err != nil {
		return false, nil
	}

	return true, nil
}
