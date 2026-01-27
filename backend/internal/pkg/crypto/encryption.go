package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"sync"
)

var (
	globalEncryptor *Encryptor
	once            sync.Once

	ErrNotInitialized = errors.New("encryptor not initialized")
	ErrInvalidKey     = errors.New("encryption key must be 32 bytes (64 hex chars)")
	ErrDecryptFailed  = errors.New("decryption failed: invalid ciphertext")
)

// Encryptor handles field-level encryption using AES-256
type Encryptor struct {
	deterministicKey []byte // For searchable fields (same input = same output)
	probabilisticKey []byte // For non-searchable fields (random IV each time)
}

// Init initializes the global encryptor with a hex-encoded master key
// The master key should be 32 bytes (64 hex characters)
func Init(masterKeyHex string) error {
	masterKey, err := hex.DecodeString(masterKeyHex)
	if err != nil || len(masterKey) != 32 {
		return ErrInvalidKey
	}

	enc, err := NewEncryptor(masterKey)
	if err != nil {
		return err
	}

	once.Do(func() {
		globalEncryptor = enc
	})

	return nil
}

// Get returns the global encryptor instance
func Get() *Encryptor {
	return globalEncryptor
}

// NewEncryptor creates a new Encryptor with derived keys from master key
func NewEncryptor(masterKey []byte) (*Encryptor, error) {
	if len(masterKey) != 32 {
		return nil, ErrInvalidKey
	}

	// Derive separate keys using HMAC-SHA256
	deterministicKey := deriveKey(masterKey, []byte("deterministic"))
	probabilisticKey := deriveKey(masterKey, []byte("probabilistic"))

	return &Encryptor{
		deterministicKey: deterministicKey,
		probabilisticKey: probabilisticKey,
	}, nil
}

// deriveKey derives a 32-byte key from master key and context using HMAC-SHA256
func deriveKey(masterKey, context []byte) []byte {
	h := hmac.New(sha256.New, masterKey)
	h.Write(context)
	return h.Sum(nil)
}

// EncryptDeterministic encrypts plaintext deterministically
// Same plaintext always produces the same ciphertext (allows equality search)
// Uses HMAC of plaintext as IV for determinism
func (e *Encryptor) EncryptDeterministic(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.deterministicKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Derive nonce from plaintext using HMAC (deterministic)
	nonce := deriveNonce(e.deterministicKey, []byte(plaintext), aesGCM.NonceSize())

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptDeterministic decrypts ciphertext encrypted with EncryptDeterministic
func (e *Encryptor) DecryptDeterministic(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrDecryptFailed
	}

	block, err := aes.NewCipher(e.deterministicKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", ErrDecryptFailed
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plaintext), nil
}

// EncryptProbabilistic encrypts plaintext with a random IV
// Same plaintext produces different ciphertext each time (more secure)
func (e *Encryptor) EncryptProbabilistic(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.probabilisticKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptProbabilistic decrypts ciphertext encrypted with EncryptProbabilistic
func (e *Encryptor) DecryptProbabilistic(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrDecryptFailed
	}

	block, err := aes.NewCipher(e.probabilisticKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", ErrDecryptFailed
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plaintext), nil
}

// EncryptNullableD encrypts a nullable string deterministically
func (e *Encryptor) EncryptNullableD(s *string) (*string, error) {
	if s == nil {
		return nil, nil
	}
	encrypted, err := e.EncryptDeterministic(*s)
	if err != nil {
		return nil, err
	}
	return &encrypted, nil
}

// DecryptNullableD decrypts a nullable string that was encrypted deterministically
func (e *Encryptor) DecryptNullableD(s *string) (*string, error) {
	if s == nil {
		return nil, nil
	}
	decrypted, err := e.DecryptDeterministic(*s)
	if err != nil {
		return nil, err
	}
	return &decrypted, nil
}

// EncryptNullableP encrypts a nullable string probabilistically
func (e *Encryptor) EncryptNullableP(s *string) (*string, error) {
	if s == nil {
		return nil, nil
	}
	encrypted, err := e.EncryptProbabilistic(*s)
	if err != nil {
		return nil, err
	}
	return &encrypted, nil
}

// DecryptNullableP decrypts a nullable string that was encrypted probabilistically
func (e *Encryptor) DecryptNullableP(s *string) (*string, error) {
	if s == nil {
		return nil, nil
	}
	decrypted, err := e.DecryptProbabilistic(*s)
	if err != nil {
		return nil, err
	}
	return &decrypted, nil
}

// deriveNonce derives a deterministic nonce from key and data using HMAC
func deriveNonce(key, data []byte, size int) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	hash := h.Sum(nil)
	return hash[:size]
}
