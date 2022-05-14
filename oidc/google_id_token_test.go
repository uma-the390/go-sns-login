package oidc

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGoogleIdTokenPayload_IsValid(t *testing.T) {
	validClientId := os.Getenv("GOOGLE_CLIENT_ID") // TODO: ちゃんとしたい

	patterns := []struct {
		desc     string
		expected error
		clientId string
		iss      string
		aud      string
		exp      int64
	}{
		{
			"valid",
			nil,
			validClientId,
			"https://accounts.google.com",
			validClientId,
			time.Now().AddDate(0, 0, 1).Unix(),
		},
		{
			"invalid iss",
			errIssMismatch,
			validClientId,
			"https://accounts.google.coms",
			validClientId,
			time.Now().AddDate(0, 0, 1).Unix(),
		},
		{
			"invalid aud",
			errAudMismatch,
			validClientId,
			"https://accounts.google.com",
			"invalid aud",
			time.Now().AddDate(0, 0, 1).Unix(),
		},
		{
			"invalid exp",
			errIdTokenExpired,
			validClientId,
			"https://accounts.google.com",
			validClientId,
			time.Now().AddDate(0, 0, -1).Unix(),
		},
	}

	for _, pattern := range patterns {
		payload := googleIdTokenPayload{
			Iss: pattern.iss,
			Aud: pattern.aud,
			Exp: pattern.exp,
		}

		actual := payload.validate(pattern.clientId)

		assert.Equal(t, pattern.expected, actual)
	}
}
