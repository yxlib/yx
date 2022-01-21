// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	ErrAesDataLenZero      = errors.New("data len is 0")
	ErrAesKeyLenZero       = errors.New("key len is 0")
	ErrAesUnpaddingTooMuch = errors.New("unpadding count more than data len")
)

// Aes encrypt.
// @param origData, data need to encrypt.
// @param key, aes encrypt key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
// @return []byte, data encrypted.
// @return error, error.
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, ErrAesDataLenZero
	}

	if len(key) == 0 {
		return nil, ErrAesKeyLenZero
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// Aes decrypt.
// @param encrypted, data need to decrypt.
// @param key, aes decrypt key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
// @return []byte, original data.
// @return error, error.
func AesDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	if len(encrypted) == 0 {
		return nil, ErrAesDataLenZero
	}

	if len(key) == 0 {
		return nil, ErrAesKeyLenZero
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData, err = PKCS7Unpadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, nil
}

// Aes encrypt with string key which encode by base64.
// @param origData, data need to encrypt.
// @param keyBase64, aes encrypt key which encode by base64.
// @return string, data encrypted which encode by base64.
// @return error, error.
func AesEncryptBase64(origData []byte, keyBase64 string) (string, error) {
	if len(origData) == 0 {
		return "", ErrAesDataLenZero
	}

	if len(keyBase64) == 0 {
		return "", ErrAesKeyLenZero
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", err
	}

	result, err := AesEncrypt(origData, key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

// Aes decrypt with string key which encode by base64.
// @param encryptedBase64, data need to decrypt which encode by base64.
// @param keyBase64, aes decrypt key which encode by base64.
// @return []byte, original data.
// @return error, error.
func AesDecryptBase64(encryptedBase64 string, keyBase64 string) ([]byte, error) {
	if len(encryptedBase64) == 0 {
		return nil, ErrAesDataLenZero
	}

	if len(keyBase64) == 0 {
		return nil, ErrAesKeyLenZero
	}

	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	b, err := AesDecrypt(encrypted, key)
	return b, err
}

func PKCS7Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	b := []byte{byte(padding)}
	paddingBytes := bytes.Repeat(b, padding)
	return append(origData, paddingBytes...)
}

func PKCS7Unpadding(decrypted []byte) ([]byte, error) {
	length := len(decrypted)
	if length == 0 {
		return nil, ErrAesDataLenZero
	}

	unpadding := int(decrypted[length-1])
	if unpadding >= length {
		return nil, ErrAesUnpaddingTooMuch
	}

	return decrypted[:(length - unpadding)], nil
}
