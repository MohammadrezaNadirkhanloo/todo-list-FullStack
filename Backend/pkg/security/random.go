package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

const otpDigits = "0123456789"

func RandomDigits(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("security: طول باید مثبت باشد، مقدار داده‌شده %d", length)
	}
	out := make([]byte, length)
	max := big.NewInt(int64(len(otpDigits)))
	for i := range out {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("security: تولید عدد تصادفی ناموفق بود: %w", err)
		}
		out[i] = otpDigits[n.Int64()]
	}
	return string(out), nil
}

func RandomBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("security: تعداد بایت باید مثبت باشد، مقدار داده‌شده %d", n)
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("security: خواندن از منبع تصادفی ناموفق بود: %w", err)
	}
	return b, nil
}

func RandomToken(byteLen int) (string, error) {
	b, err := RandomBytes(byteLen)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func ConstantTimeEquals(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}