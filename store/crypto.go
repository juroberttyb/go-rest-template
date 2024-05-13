package store

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/A-pen-app/logging"
)

type keyStore struct {
	c *kms.KeyManagementClient
}

func NewCrypto(ctx context.Context) Crypto {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		panic(err)
	}

	return &keyStore{
		c: client,
	}
}

func (k *keyStore) Decrypt(ctx context.Context, keyID string, base64Ciphertext string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req := &kmspb.AsymmetricDecryptRequest{
		Name:       keyID,
		Ciphertext: ciphertext,
	}

	result, err := k.c.AsymmetricDecrypt(ctx, req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result.Plaintext, nil
}

func (k *keyStore) CreateKey(ctx context.Context, keyRing, keyName string) error {
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      keyRing,
		CryptoKeyId: keyName,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA512,
			},
		},
	}

	_, err := k.c.CreateCryptoKey(ctx, req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (k *keyStore) Sign(ctx context.Context, keyID string, msg string) (string, error) {

	plaintext := []byte(msg)
	digest := sha256.New()
	if _, err := digest.Write(plaintext); err != nil {
		logging.Errorw(ctx, "calculate SHA-256 hash for input msg failed", "err", err)
		return "", err
	}

	req := &kmspb.AsymmetricSignRequest{
		Name: keyID,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest.Sum(nil),
			},
		},
	}

	result, err := k.c.AsymmetricSign(ctx, req)
	if err != nil {
		logging.Errorw(ctx, "create signature failed", "err", err)
		return "", err
	}
	sig := base64.URLEncoding.EncodeToString(result.Signature)
	return sig, nil
}

func (k *keyStore) Verify(ctx context.Context, keyID, msg, signature string) error {
	result, err := k.c.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: keyID,
	})
	if err != nil {
		logging.Errorw(ctx, "store get public key failed", "err", err, "keyID", keyID)
		return err
	}

	block, _ := pem.Decode([]byte(result.Pem))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logging.Errorw(ctx, "store parse public key failed", "err", err, "keyID", keyID)
		return err
	}
	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		logging.Errorw(ctx, "public key format is not rsa", "keyID", keyID)
		return err
	}

	// Verify the RSA signature.
	rawSig, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		logging.Errorw(ctx, "signature format is not base64", "signature", signature, "err", err)
		return err
	}
	plaintext := []byte(msg)
	d := sha256.New()
	if _, err := d.Write(plaintext); err != nil {
		logging.Errorw(ctx, "calculate SHA-256 hash for input msg failed", "err", err)
		return err
	}
	digest := d.Sum(nil)
	if err := rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, digest[:], rawSig); err != nil {
		logging.Errorw(ctx, "store verify signature failed", "msg", msg, "keyID", keyID, "err", err, "sig", signature)
		return err
	}

	return nil
}

func (k *keyStore) GetPublicKey(ctx context.Context, keyID string) (string, error) {
	req := &kmspb.GetPublicKeyRequest{
		Name: keyID,
	}
	result, err := k.c.GetPublicKey(ctx, req)
	if err != nil {
		logging.Errorw(ctx, "store get public key failed", "err", err, "keyID", keyID)
		return "", err
	}
	key := base64.StdEncoding.EncodeToString([]byte(result.Pem))
	return key, nil
}
