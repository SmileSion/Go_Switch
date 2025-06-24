package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
)

var secretKey []byte

func init() {
	var err error
	secretKey, err = hex.DecodeString("dd712682983e1d71b80bdf2f7a15353846d48552b886242dbab3d0fbafb95ec1")
	if err != nil {
		log.Fatalf("密钥解码失败: %v", err)
	}
}

func Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return string(data), nil
}

// 可选：加密函数
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, secretKey[:aes.BlockSize])
	stream := cipher.NewCFBEncrypter(block, iv)
	data := []byte(plainText)
	stream.XORKeyStream(data, data)
	result := append(iv, data...)
	return base64.StdEncoding.EncodeToString(result), nil
}
