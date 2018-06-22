package server

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher hashes passwords and verifies passwords against the generated hashes.
type PasswordHasher interface {
	HashAndSalt(plainPassword string) string
	VerifyPassword(hashedPassword, plainPassword string) bool
}

// DefaultPasswordHasher implements PasswordHasher.
type DefaultPasswordHasher struct{}

// HashAndSalt one-way hashes the password.
func (p *DefaultPasswordHasher) HashAndSalt(plainPassword string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// VerifyPassword verifies if the plainPassword matches the hashedPassword.
func (p *DefaultPasswordHasher) VerifyPassword(hashedPassword, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	if err != nil {
		return false
	}
	return true
}
