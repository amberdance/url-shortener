package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/amberdance/url-shortener/internal/domain/errs"
)

type CookieAuth struct {
	secret []byte
}

func NewCookieAuth(secret string) *CookieAuth {
	return &CookieAuth{secret: []byte(secret)}
}

func (a *CookieAuth) Sign(userID string) string {
	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(userID))
	sign := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	payload := base64.RawURLEncoding.EncodeToString([]byte(userID))

	return payload + "." + sign
}

func (a *CookieAuth) Verify(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", errs.InvalidArgumentError("invalid token format")
	}

	payload, sig := parts[0], parts[1]
	userBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", errs.InvalidArgumentError("invalid base64 payload")
	}

	userID := string(userBytes)

	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(userID))
	expectedSig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", errs.UnauthorizedError("invalid signature")
	}

	return userID, nil
}
