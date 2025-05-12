package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"zhaoxin2025/config"
)

// Encrypt 加密函数
func Encrypt(text string) (string, error) {
	key := []byte(config.Config.AppSalt)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data, prepending the nonce to the ciphertext
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密函数
func Decrypt(ciphertext string) (string, error) {
	key := []byte(config.Config.AppSalt)
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	nonce, encryptedMessage := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", err // Decryption failed (possibly due to tampered ciphertext or wrong key)
	}

	return string(plaintext), nil
}
