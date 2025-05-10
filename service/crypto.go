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
	// 创建 AES 加密器
	block, err := aes.NewCipher([]byte(config.Config.AppSalt))
	if err != nil {
		return "", err
	}

	// 生成随机 IV（初始化向量）
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 使用 CFB 模式加密
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// 返回 Base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密函数
func Decrypt(ciphertext string) (string, error) {
	// 创建 AES 加密器
	block, err := aes.NewCipher([]byte(config.Config.AppSalt))
	if err != nil {
		return "", err
	}

	// 解码 Base64 密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 提取 IV（初始化向量）
	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("密文长度不足")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	// 使用 CFB 模式解密
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
