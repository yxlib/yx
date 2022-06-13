// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrRsaDataLenZero = errors.New("data len is 0")
	ErrRsaKeyLenZero  = errors.New("key len is 0")
)

func RsaGenerateKey(bits int) (pubKey []byte, priKey []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	pubKey = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	priKey = x509.MarshalPKCS1PrivateKey(key)
	return pubKey, priKey, nil
}

func RsaGenerateKeyPem(bits int) (pemPubKey string, pemPriKey string, err error) {
	pubKey, priKey, err := RsaGenerateKey(bits)
	if err != nil {
		return "", "", err
	}

	pubKeyBlock := &pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	pemPubKey = string(pem.EncodeToMemory(pubKeyBlock))

	priKeyBlock := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priKey,
	}

	pemPriKey = string(pem.EncodeToMemory(priKeyBlock))
	return pemPubKey, pemPriKey, nil
}

func RsaEncrypt(origData []byte, pubKey []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(pubKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	publicKey, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
	return encrypted, err
}

func RsaDecrypt(encrypted []byte, priKey []byte) ([]byte, error) {
	if len(encrypted) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(priKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}

	oriData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	return oriData, err
}

func RsaEncryptPem(origData []byte, pemPubKey string) ([]byte, error) {
	if len(origData) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(pemPubKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPubKey))
	return RsaEncrypt(origData, block.Bytes)
}

func RsaDecryptPem(encrypted []byte, pemPriKey string) ([]byte, error) {
	if len(encrypted) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(pemPriKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPriKey))
	return RsaDecrypt(encrypted, block.Bytes)
}

func RsaSign(origData []byte, priKey []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(priKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}

	hashed := sha512.Sum512(origData)
	signData, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashed[:])
	return signData, err
}

func RsaVerify(origData []byte, signData []byte, pubKey []byte) (bool, error) {
	if len(origData) == 0 {
		return false, ErrRsaDataLenZero
	}

	if len(signData) == 0 {
		return false, ErrRsaDataLenZero
	}

	if len(pubKey) == 0 {
		return false, ErrRsaKeyLenZero
	}

	publicKey, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return false, err
	}

	hashed := sha512.Sum512(origData)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashed[:], signData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func RsaSignPem(origData []byte, pemPriKey string) ([]byte, error) {
	if len(origData) == 0 {
		return nil, ErrRsaDataLenZero
	}

	if len(pemPriKey) == 0 {
		return nil, ErrRsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPriKey))
	return RsaSign(origData, block.Bytes)
}

func RsaVerifyPem(origData []byte, signData []byte, pemPubKey string) (bool, error) {
	if len(origData) == 0 {
		return false, ErrRsaDataLenZero
	}

	if len(signData) == 0 {
		return false, ErrRsaDataLenZero
	}

	if len(pemPubKey) == 0 {
		return false, ErrRsaKeyLenZero
	}

	block, _ := pem.Decode([]byte(pemPubKey))
	return RsaVerify(origData, signData, block.Bytes)
}
