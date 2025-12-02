package crypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

// GenerateKeyPair generates RSA public/private key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// PrivateKeyToString converts private key to PEM string
func PrivateKeyToString(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// PublicKeyToString converts public key to PEM string
func PublicKeyToString(publicKey *rsa.PublicKey) string {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

// StringToPrivateKey converts PEM string to private key
func StringToPrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// StringToPublicKey converts PEM string to public key
func StringToPublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return publicKey, nil
}

// GenerateWalletID generates a unique wallet ID from public key
func GenerateWalletID(publicKey *rsa.PublicKey) string {
	publicKeyStr := PublicKeyToString(publicKey)
	hash := sha256.Sum256([]byte(publicKeyStr))
	return hex.EncodeToString(hash[:])
}

// SignData signs data with private key
func SignData(data string, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a signature with public key
func VerifySignature(data string, signature string, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256([]byte(data))
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
}

// EncryptPrivateKey encrypts private key using AES
func EncryptPrivateKey(privateKeyStr string) (string, error) {
	aesKeyStr := os.Getenv("AES_ENCRYPTION_KEY")
	aesKey, err := base64.StdEncoding.DecodeString(aesKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode AES key: %v", err)
	}

	if len(aesKey) != 32 {
		return "", fmt.Errorf("AES key must be 32 bytes, got %d bytes", len(aesKey))
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(privateKeyStr), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey decrypts private key using AES
func DecryptPrivateKey(encryptedKey string) (string, error) {
	aesKeyStr := os.Getenv("AES_ENCRYPTION_KEY")
	aesKey, err := base64.StdEncoding.DecodeString(aesKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode AES key: %v", err)
	}

	if len(aesKey) != 32 {
		return "", fmt.Errorf("AES key must be 32 bytes, got %d bytes", len(aesKey))
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashSHA256 computes SHA-256 hash of data
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CreateTransactionPayload creates the payload to be signed for a transaction
func CreateTransactionPayload(senderID, receiverID string, amount float64, timestamp, note string) string {
	return fmt.Sprintf("%s%s%.8f%s%s", senderID, receiverID, amount, timestamp, note)
}
