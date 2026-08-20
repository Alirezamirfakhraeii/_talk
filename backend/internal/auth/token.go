package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateToken() (string, string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", fmt.Errorf(
			"generate secure random token: %w",
			err,
		)
	}

	token := base64.RawURLEncoding.EncodeToString(
		randomBytes,
	)

	tokenHash := HashToken(token)

	return token, tokenHash, nil
}

func HashToken(
	token string,
) string {
	hash := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		hash[:],
	)
}
