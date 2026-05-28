package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hash, err := HashPassword("mypassword")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if hash == "" {
			t.Fatal("expected non-empty hash")
		}
		if hash == "mypassword" {
			t.Fatal("hash should not equal plaintext password")
		}
	})

	t.Run("different_passwords_different_hashes", func(t *testing.T) {
		h1, _ := HashPassword("password1")
		h2, _ := HashPassword("password2")
		if h1 == h2 {
			t.Fatal("different passwords should produce different hashes")
		}
	})

	t.Run("same_password_different_hashes", func(t *testing.T) {
		h1, _ := HashPassword("password")
		h2, _ := HashPassword("password")
		if h1 == h2 {
			t.Fatal("same password should produce different hashes due to salt")
		}
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("correct_password", func(t *testing.T) {
		hash, _ := HashPassword("correct")
		if !CheckPassword("correct", hash) {
			t.Fatal("expected password to match")
		}
	})

	t.Run("wrong_password", func(t *testing.T) {
		hash, _ := HashPassword("correct")
		if CheckPassword("wrong", hash) {
			t.Fatal("expected password NOT to match")
		}
	})

	t.Run("empty_password", func(t *testing.T) {
		hash, _ := HashPassword("")
		if !CheckPassword("", hash) {
			t.Fatal("empty password should match")
		}
	})
}
