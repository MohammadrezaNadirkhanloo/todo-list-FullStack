package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretLength = 32

var (
	ErrTokenExpired = errors.New("security: توکن منقضی شده است")
	ErrTokenInvalid = errors.New("security: توکن نامعتبر است")
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int64    `json:"uid"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func (c *Claims) HasRole(allowed ...string) bool {
	for _, role := range c.Roles {
		for _, want := range allowed {
			if role == want {
				return true
			}
		}
	}
	return false
}

type TokenManager struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
}

func NewTokenManager(secret, issuer, audience string, ttl time.Duration) (*TokenManager, error) {
	if len(secret) < minSecretLength {
		return nil, fmt.Errorf(
			"security: کلید JWT باید حداقل %d بایت باشد (طول فعلی: %d). برای تولید: openssl rand -base64 48",
			minSecretLength, len(secret),
		)
	}
	if ttl <= 0 {
		return nil, errors.New("security: مدت اعتبار توکن باید مثبت باشد")
	}
	return &TokenManager{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
		ttl:      ttl,
	}, nil
}

func (m *TokenManager) TTL() time.Duration { return m.ttl }

type IssuedToken struct {
	Token     string
	ID        string
	ExpiresAt time.Time
}

func (m *TokenManager) Issue(userID int64, username string, roles []string) (IssuedToken, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	jti, err := RandomToken(16)
	if err != nil {
		return IssuedToken{}, err
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings{m.audience},
			Subject:   fmt.Sprintf("%d", userID),
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:   userID,
		Username: username,
		Roles:    roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return IssuedToken{}, fmt.Errorf("security: امضای توکن ناموفق بود: %w", err)
	}
	return IssuedToken{Token: signed, ID: jti, ExpiresAt: expiresAt}, nil
}

func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(_ *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	return claims, nil
}