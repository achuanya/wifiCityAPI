package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// EncryptedData 结构用于封装加密后的数据
type EncryptedData struct {
	Data string `json:"data"` // Base64 编码的密文
	IV   string `json:"iv"`   // Base64 编码的初始化向量
	Tag  string `json:"tag"`  // Base64 编码的认证标签
}

// Encrypt 使用 AES-256-GCM 加密数据
// key 必须是 16, 24, 或 32 字节
func Encrypt(plaintext []byte, key []byte) (*EncryptedData, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// GCM 的 Seal 方法返回的 ciphertext 已经包含了 tag，我们需要分离它们
	tagSize := gcm.Overhead()
	realCiphertext := ciphertext[:len(ciphertext)-tagSize]
	tag := ciphertext[len(ciphertext)-tagSize:]

	return &EncryptedData{
		Data: base64.StdEncoding.EncodeToString(realCiphertext),
		IV:   base64.StdEncoding.EncodeToString(nonce),
		Tag:  base64.StdEncoding.EncodeToString(tag),
	}, nil
}

// Decrypt 使用 AES-256-GCM 解密数据
func Decrypt(encryptedData *EncryptedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedData.IV)
	if err != nil {
		return nil, fmt.Errorf("iV 解码失败: %w", err)
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("无效的 IV 长度")
	}

	tag, err := base64.StdEncoding.DecodeString(encryptedData.Tag)
	if err != nil {
		return nil, fmt.Errorf("tag 解码失败: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData.Data)
	if err != nil {
		return nil, fmt.Errorf("密文解码失败: %w", err)
	}

	// 将密文和tag合并回gcm Open需要的格式
	fullCiphertext := append(ciphertext, tag...)

	plaintext, err := gcm.Open(nil, nonce, fullCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %w", err)
	}

	return plaintext, nil
}

// GenerateSignature 使用 HMAC-SHA256 生成签名
func GenerateSignature(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// ValidateSignature 验证签名
func ValidateSignature(message, signature string, key []byte) bool {
	expectedSignature := GenerateSignature(message, key)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
