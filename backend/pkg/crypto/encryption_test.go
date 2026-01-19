package crypto

import (
	"testing"
)

const testMasterKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestNewEncryptor(t *testing.T) {
	masterKey := make([]byte, 32)
	enc, err := NewEncryptor(masterKey)
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}
	if enc == nil {
		t.Fatal("NewEncryptor returned nil")
	}
}

func TestNewEncryptorInvalidKey(t *testing.T) {
	// Too short
	_, err := NewEncryptor(make([]byte, 16))
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey for short key, got %v", err)
	}

	// Too long
	_, err = NewEncryptor(make([]byte, 64))
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey for long key, got %v", err)
	}
}

func TestInit(t *testing.T) {
	// Note: Init uses sync.Once, so this test should run in isolation
	// For now, we test that it doesn't panic with valid input
	err := Init(testMasterKey)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	enc := Get()
	if enc == nil {
		t.Fatal("Get() returned nil after Init")
	}
}

func TestInitInvalidKey(t *testing.T) {
	err := Init("invalid")
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey for invalid hex, got %v", err)
	}

	err = Init("0123456789abcdef") // Too short
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey for short key, got %v", err)
	}
}

func TestDeterministicEncryption(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	tests := []struct {
		name      string
		plaintext string
	}{
		{"simple", "hello world"},
		{"email", "test@example.com"},
		{"phone", "+6281234567890"},
		{"unicode", "日本語テスト"},
		{"special chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := enc.EncryptDeterministic(tt.plaintext)
			if err != nil {
				t.Fatalf("EncryptDeterministic failed: %v", err)
			}

			// Ciphertext should be different from plaintext
			if ciphertext == tt.plaintext {
				t.Error("Ciphertext equals plaintext")
			}

			// Decrypt
			decrypted, err := enc.DecryptDeterministic(ciphertext)
			if err != nil {
				t.Fatalf("DecryptDeterministic failed: %v", err)
			}

			// Should match original
			if decrypted != tt.plaintext {
				t.Errorf("Decrypted text doesn't match: got %q, want %q", decrypted, tt.plaintext)
			}

			// Same plaintext should produce same ciphertext (deterministic)
			ciphertext2, _ := enc.EncryptDeterministic(tt.plaintext)
			if ciphertext2 != ciphertext {
				t.Error("Deterministic encryption produced different ciphertext for same plaintext")
			}
		})
	}
}

func TestDeterministicEncryptionEmptyString(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	ciphertext, err := enc.EncryptDeterministic("")
	if err != nil {
		t.Fatalf("EncryptDeterministic empty string failed: %v", err)
	}
	if ciphertext != "" {
		t.Errorf("Expected empty ciphertext for empty plaintext, got %q", ciphertext)
	}

	decrypted, err := enc.DecryptDeterministic("")
	if err != nil {
		t.Fatalf("DecryptDeterministic empty string failed: %v", err)
	}
	if decrypted != "" {
		t.Errorf("Expected empty decrypted for empty ciphertext, got %q", decrypted)
	}
}

func TestProbabilisticEncryption(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	tests := []struct {
		name      string
		plaintext string
	}{
		{"simple", "hello world"},
		{"name", "John Doe"},
		{"address", "123 Main St, City, Country"},
		{"unicode", "日本語テスト"},
		{"long text", "This is a longer piece of text that might be stored in a remarks field."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := enc.EncryptProbabilistic(tt.plaintext)
			if err != nil {
				t.Fatalf("EncryptProbabilistic failed: %v", err)
			}

			// Ciphertext should be different from plaintext
			if ciphertext == tt.plaintext {
				t.Error("Ciphertext equals plaintext")
			}

			// Decrypt
			decrypted, err := enc.DecryptProbabilistic(ciphertext)
			if err != nil {
				t.Fatalf("DecryptProbabilistic failed: %v", err)
			}

			// Should match original
			if decrypted != tt.plaintext {
				t.Errorf("Decrypted text doesn't match: got %q, want %q", decrypted, tt.plaintext)
			}

			// Same plaintext should produce DIFFERENT ciphertext (probabilistic)
			ciphertext2, _ := enc.EncryptProbabilistic(tt.plaintext)
			if ciphertext2 == ciphertext {
				t.Error("Probabilistic encryption produced same ciphertext for same plaintext")
			}

			// But both should decrypt to same plaintext
			decrypted2, _ := enc.DecryptProbabilistic(ciphertext2)
			if decrypted2 != tt.plaintext {
				t.Errorf("Second decryption doesn't match: got %q, want %q", decrypted2, tt.plaintext)
			}
		})
	}
}

func TestProbabilisticEncryptionEmptyString(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	ciphertext, err := enc.EncryptProbabilistic("")
	if err != nil {
		t.Fatalf("EncryptProbabilistic empty string failed: %v", err)
	}
	if ciphertext != "" {
		t.Errorf("Expected empty ciphertext for empty plaintext, got %q", ciphertext)
	}

	decrypted, err := enc.DecryptProbabilistic("")
	if err != nil {
		t.Fatalf("DecryptProbabilistic empty string failed: %v", err)
	}
	if decrypted != "" {
		t.Errorf("Expected empty decrypted for empty ciphertext, got %q", decrypted)
	}
}

func TestNullableEncryption(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	// Test nil input
	resultD, err := enc.EncryptNullableD(nil)
	if err != nil {
		t.Fatalf("EncryptNullableD nil failed: %v", err)
	}
	if resultD != nil {
		t.Error("Expected nil result for nil input (deterministic)")
	}

	resultP, err := enc.EncryptNullableP(nil)
	if err != nil {
		t.Fatalf("EncryptNullableP nil failed: %v", err)
	}
	if resultP != nil {
		t.Error("Expected nil result for nil input (probabilistic)")
	}

	// Test non-nil input
	input := "test value"
	encryptedD, err := enc.EncryptNullableD(&input)
	if err != nil {
		t.Fatalf("EncryptNullableD failed: %v", err)
	}
	if encryptedD == nil {
		t.Fatal("Expected non-nil result for non-nil input (deterministic)")
	}

	decryptedD, err := enc.DecryptNullableD(encryptedD)
	if err != nil {
		t.Fatalf("DecryptNullableD failed: %v", err)
	}
	if *decryptedD != input {
		t.Errorf("Deterministic nullable roundtrip failed: got %q, want %q", *decryptedD, input)
	}

	encryptedP, err := enc.EncryptNullableP(&input)
	if err != nil {
		t.Fatalf("EncryptNullableP failed: %v", err)
	}
	if encryptedP == nil {
		t.Fatal("Expected non-nil result for non-nil input (probabilistic)")
	}

	decryptedP, err := enc.DecryptNullableP(encryptedP)
	if err != nil {
		t.Fatalf("DecryptNullableP failed: %v", err)
	}
	if *decryptedP != input {
		t.Errorf("Probabilistic nullable roundtrip failed: got %q, want %q", *decryptedP, input)
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	enc, _ := NewEncryptor(make([]byte, 32))

	// Invalid base64
	_, err := enc.DecryptDeterministic("not-valid-base64!!!")
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed for invalid base64, got %v", err)
	}

	// Valid base64 but invalid ciphertext (too short)
	_, err = enc.DecryptDeterministic("dGVzdA==") // "test" in base64
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed for short ciphertext, got %v", err)
	}

	// Valid base64 but tampered ciphertext
	ciphertext, _ := enc.EncryptDeterministic("hello")
	tampered := ciphertext[:len(ciphertext)-2] + "XX"
	_, err = enc.DecryptDeterministic(tampered)
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed for tampered ciphertext, got %v", err)
	}
}

func TestDifferentKeysProduceDifferentCiphertext(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	key2[0] = 1 // Different key

	enc1, _ := NewEncryptor(key1)
	enc2, _ := NewEncryptor(key2)

	plaintext := "test message"

	cipher1, _ := enc1.EncryptDeterministic(plaintext)
	cipher2, _ := enc2.EncryptDeterministic(plaintext)

	if cipher1 == cipher2 {
		t.Error("Different keys produced same ciphertext")
	}

	// Decrypting with wrong key should fail
	_, err := enc2.DecryptDeterministic(cipher1)
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed when decrypting with wrong key, got %v", err)
	}
}

func BenchmarkEncryptDeterministic(b *testing.B) {
	enc, _ := NewEncryptor(make([]byte, 32))
	plaintext := "test@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.EncryptDeterministic(plaintext)
	}
}

func BenchmarkDecryptDeterministic(b *testing.B) {
	enc, _ := NewEncryptor(make([]byte, 32))
	ciphertext, _ := enc.EncryptDeterministic("test@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.DecryptDeterministic(ciphertext)
	}
}

func BenchmarkEncryptProbabilistic(b *testing.B) {
	enc, _ := NewEncryptor(make([]byte, 32))
	plaintext := "This is a longer text that might be in a remarks field"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.EncryptProbabilistic(plaintext)
	}
}

func BenchmarkDecryptProbabilistic(b *testing.B) {
	enc, _ := NewEncryptor(make([]byte, 32))
	ciphertext, _ := enc.EncryptProbabilistic("This is a longer text that might be in a remarks field")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.DecryptProbabilistic(ciphertext)
	}
}
