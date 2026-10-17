package security

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

var ErrMismatch = errors.New("security: رمز عبور مطابقت ندارد")

var ErrUnknownHashFormat = errors.New("security: قالب هش ناشناخته است")

type Hasher interface {
	Hash(plain string) (string, error)
	Verify(plain, encoded string) (needsRehash bool, err error)
}

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DefaultArgon2Params() Argon2Params {
	parallelism := uint8(2)
	if n := runtime.NumCPU(); n < 2 {
		parallelism = 1
	}
	return Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: parallelism,
		SaltLength:  16,
		KeyLength:   32,
	}
}

type Argon2Hasher struct {
	params Argon2Params
}

func NewArgon2Hasher(params Argon2Params) *Argon2Hasher {
	if params.Memory == 0 {
		params = DefaultArgon2Params()
	}
	return &Argon2Hasher{params: params}
}

func (h *Argon2Hasher) Hash(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("security: رمز عبور نمی‌تواند خالی باشد")
	}

	salt, err := RandomBytes(int(h.params.SaltLength))
	if err != nil {
		return "", err
	}

	key := argon2.IDKey(
		[]byte(plain),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.params.Memory,
		h.params.Iterations,
		h.params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func (h *Argon2Hasher) Verify(plain, encoded string) (bool, error) {
	switch {
	case strings.HasPrefix(encoded, "$argon2id$"):
		return h.verifyArgon2(plain, encoded)
	case strings.HasPrefix(encoded, "$2a$"),
		strings.HasPrefix(encoded, "$2b$"),
		strings.HasPrefix(encoded, "$2y$"):
		if err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(plain)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return false, ErrMismatch
			}
			return false, fmt.Errorf("security: بررسی هش bcrypt ناموفق بود: %w", err)
		}
		return true, nil
	default:
		return false, ErrUnknownHashFormat
	}
}

func (h *Argon2Hasher) verifyArgon2(plain, encoded string) (bool, error) {
	params, salt, key, err := decodeArgon2Hash(encoded)
	if err != nil {
		return false, err
	}

	candidate := argon2.IDKey(
		[]byte(plain),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		uint32(len(key)),
	)

	if subtle.ConstantTimeCompare(key, candidate) != 1 {
		return false, ErrMismatch
	}

	needsRehash := params.Memory < h.params.Memory ||
		params.Iterations < h.params.Iterations ||
		params.Parallelism < h.params.Parallelism

	return needsRehash, nil
}

func decodeArgon2Hash(encoded string) (Argon2Params, []byte, []byte, error) {
	var p Argon2Params

	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return p, nil, nil, ErrUnknownHashFormat
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return p, nil, nil, ErrUnknownHashFormat
	}
	if version != argon2.Version {
		return p, nil, nil, fmt.Errorf("security: نسخه‌ی ناسازگار argon2: %d", version)
	}

	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism); err != nil {
		return p, nil, nil, ErrUnknownHashFormat
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return p, nil, nil, ErrUnknownHashFormat
	}
	key, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return p, nil, nil, ErrUnknownHashFormat
	}

	p.SaltLength = uint32(len(salt))
	p.KeyLength = uint32(len(key))
	return p, salt, key, nil
}

func (h *Argon2Hasher) DummyVerify() {
	salt := make([]byte, h.params.SaltLength)
	_ = argon2.IDKey(
		[]byte("dummy-password-for-constant-time"),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)
}