package yx

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
)

var (
	ErrEcdsaDataLenZero     = errors.New("data len is 0")
	ErrEcdsaKeyLenZero      = errors.New("key len is 0")
	ErrEcdsaParsePubKeyFail = errors.New("parse public key failed")
)

type ElliptcType int

const (
	ELLIPTIC_TYPE_P521 ElliptcType = iota
	ELLIPTIC_TYPE_P384
	ELLIPTIC_TYPE_P256
	ELLIPTIC_TYPE_P224
)

func EcdsaGenerateKey(ellType ElliptcType) (pubKey []byte, priKey []byte, err error) {
	var curve elliptic.Curve
	if ellType == ELLIPTIC_TYPE_P521 {
		curve = elliptic.P521()
	} else if ellType == ELLIPTIC_TYPE_P384 {
		curve = elliptic.P384()
	} else if ellType == ELLIPTIC_TYPE_P256 {
		curve = elliptic.P256()
	} else if ellType == ELLIPTIC_TYPE_P224 {
		curve = elliptic.P224()
	}

	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	priKey, err = x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	return pubKey, priKey, nil
}

func EcdsaGenerateKeyPem(ellType ElliptcType) (pemPubKey string, pemPriKey string, err error) {
	pubKey, priKey, err := EcdsaGenerateKey(ellType)
	if err != nil {
		return "", "", err
	}

	pubKeyBlock := &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	pemPubKey = string(pem.EncodeToMemory(pubKeyBlock))

	priKeyBlock := &pem.Block{
		Type:    "EC PRIVATE KEY",
		Headers: nil,
		Bytes:   priKey,
	}

	pemPriKey = string(pem.EncodeToMemory(priKeyBlock))
	return pemPubKey, pemPriKey, nil
}

func EcdsaSign(origData []byte, priKey []byte) ([]byte, []byte, error) {
	if len(origData) == 0 {
		return nil, nil, ErrEcdsaDataLenZero
	}

	if len(priKey) == 0 {
		return nil, nil, ErrEcdsaKeyLenZero
	}

	privateKey, err := x509.ParseECPrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}

	hashed := sha512.Sum512(origData)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])
	if err != nil {
		return nil, nil, err
	}

	rText, err := r.MarshalText()
	if err != nil {
		return nil, nil, err
	}

	sText, err := s.MarshalText()
	if err != nil {
		return nil, nil, err
	}

	return rText, sText, err
}

func EcdsaVerify(origData []byte, rText []byte, sText []byte, pubKey []byte) (bool, error) {
	if len(origData) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(rText) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(sText) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(pubKey) == 0 {
		return false, ErrEcdsaKeyLenZero
	}

	key, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return false, err
	}

	publicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return false, ErrEcdsaParsePubKeyFail
	}

	hashed := sha512.Sum512(origData)
	var r, s big.Int
	err = r.UnmarshalText(rText)
	if err != nil {
		return false, err
	}

	err = s.UnmarshalText(sText)
	if err != nil {
		return false, err
	}

	bSucc := ecdsa.Verify(publicKey, hashed[:], &r, &s)
	return bSucc, nil
}

func EcdsaSignPem(origData []byte, pemPriKey string) ([]byte, []byte, error) {
	if len(origData) == 0 {
		return nil, nil, ErrEcdsaDataLenZero
	}

	if len(pemPriKey) == 0 {
		return nil, nil, ErrEcdsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPriKey))
	return EcdsaSign(origData, block.Bytes)
}

func EcdsaVerifyPem(origData []byte, rText []byte, sText []byte, pemPubKey string) (bool, error) {
	if len(origData) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(rText) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(sText) == 0 {
		return false, ErrEcdsaDataLenZero
	}

	if len(pemPubKey) == 0 {
		return false, ErrEcdsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPubKey))
	return EcdsaVerify(origData, rText, sText, block.Bytes)
}
