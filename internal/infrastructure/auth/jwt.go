package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gsystes/backend/internal/infrastructure/config"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    RoleID   uint   `json:"role_id"`
    jwt.RegisteredClaims
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

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
    cfg := config.GetConfig().JWT
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(cfg.Secret), nil
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