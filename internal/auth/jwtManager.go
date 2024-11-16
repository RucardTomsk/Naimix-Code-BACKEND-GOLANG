package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type JWTManager struct {
	signingKey string
	timeToLive time.Duration
}

func NewJWTManager(signingKey string, timeToLive time.Duration) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{signingKey: signingKey, timeToLive: timeToLive}, nil
}

func (m *JWTManager) NewJWT(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.timeToLive).Unix(),
		Subject:   userId.String(),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}
