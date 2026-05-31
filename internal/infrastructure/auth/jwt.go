package auth

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gsystes/backend/internal/infrastructure/config"
)

type Claims struct {
	UserID   uint                   `json:"user_id"`
	Username string                 `json:"username"`
	RoleID   uint                   `json:"role_id"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
	jwt.RegisteredClaims
}

type TokenService interface {
	GenerateToken(userID uint, username string, roleID uint) (string, error)
}

type defaultTokenService struct{}

func NewTokenService() TokenService {
	return &defaultTokenService{}
}

func (s *defaultTokenService) GenerateToken(userID uint, username string, roleID uint) (string, error) {
	return GenerateToken(userID, username, roleID)
}

func getSigningMethod(method string) jwt.SigningMethod {
	switch method {
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES384":
		return jwt.SigningMethodES384
	case "ES512":
		return jwt.SigningMethodES512
	default:
		return jwt.SigningMethodHS256
	}
}

func getSigningKey(cfg config.JWTConfig) (interface{}, error) {
	method := getSigningMethod(cfg.SigningMethod)
	switch method.(type) {
	case *jwt.SigningMethodHMAC:
		return []byte(cfg.Secret), nil
	case *jwt.SigningMethodRSA:
		return loadRSAPrivateKey()
	case *jwt.SigningMethodECDSA:
		return loadECDSAPrivateKey()
	default:
		return []byte(cfg.Secret), nil
	}
}

func loadRSAPrivateKey() (*rsa.PrivateKey, error) {
	keyPath := os.Getenv("GSYSTES_JWT_RSA_PRIVATE_KEY")
	if keyPath == "" {
		keyPath = "config/jwt_rsa_private.pem"
	}
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(data)
}

func loadECDSAPrivateKey() (interface{}, error) {
	keyPath := os.Getenv("GSYSTES_JWT_ECDSA_PRIVATE_KEY")
	if keyPath == "" {
		keyPath = "config/jwt_ecdsa_private.pem"
	}
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseECPrivateKeyFromPEM(data)
}

func GenerateToken(userID uint, username string, roleID uint) (string, error) {
	cfg := config.GetConfig().JWT
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
		},
	}

	method := getSigningMethod(cfg.SigningMethod)
	token := jwt.NewWithClaims(method, claims)
	key, err := getSigningKey(cfg)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig().JWT
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		method := getSigningMethod(cfg.SigningMethod)
		if token.Method.Alg() != method.Alg() {
			return nil, errors.New("unexpected signing method: " + token.Method.Alg())
		}
		return getSigningKey(cfg)
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
