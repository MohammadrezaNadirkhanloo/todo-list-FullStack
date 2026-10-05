package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

const csrfKeyLabel = "goclean/csrf-token/v1"

const csrfNonceBytes = 32

const maxCSRFTokenLen = 256

var ErrCSRFInvalid = errors.New("security: توکن CSRF نامعتبر است")

type CSRFSigner struct {
	key []byte
}

func NewCSRFSigner(accessSecret string) (*CSRFSigner, error) {
	if len(accessSecret) < minSecretLength {
		return nil, fmt.Errorf(
			"security: کلید پایه‌ی CSRF باید حداقل %d بایت باشد (طول فعلی: %d)",
			minSecretLength, len(accessSecret),
		)
	}
	mac := hmac.New(sha256.New, []byte(accessSecret))
	mac.Write([]byte(csrfKeyLabel))
	return &CSRFSigner{key: mac.Sum(nil)}, nil
}

func (s *CSRFSigner) Issue(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("security: شناسه‌ی نشست برای توکن CSRF خالی است")
	}
	nonce, err := RandomBytes(csrfNonceBytes)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	return enc.EncodeToString(nonce) + "." + enc.EncodeToString(s.sign(sessionID, nonce)), nil
}

func (s *CSRFSigner) Verify(token, sessionID string) error {
	if sessionID == "" || token == "" || len(token) > maxCSRFTokenLen {
		return ErrCSRFInvalid
	}

	noncePart, sigPart, ok := strings.Cut(token, ".")
	if !ok {
		return ErrCSRFInvalid
	}

	enc := base64.RawURLEncoding
	nonce, err := enc.DecodeString(noncePart)
	if err != nil || len(nonce) != csrfNonceBytes {
		return ErrCSRFInvalid
	}
	sig, err := enc.DecodeString(sigPart)
	if err != nil {
		return ErrCSRFInvalid
	}

	if !hmac.Equal(sig, s.sign(sessionID, nonce)) {
		return ErrCSRFInvalid
	}
	return nil
}

func (s *CSRFSigner) sign(sessionID string, nonce []byte) []byte {
	mac := hmac.New(sha256.New, s.key)
	var n [8]byte
	binary.BigEndian.PutUint64(n[:], uint64(len(sessionID)))
	mac.Write(n[:])
	mac.Write([]byte(sessionID))
	binary.BigEndian.PutUint64(n[:], uint64(len(nonce)))
	mac.Write(n[:])
	mac.Write(nonce)
	return mac.Sum(nil)
}