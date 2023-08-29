package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type JWTManager interface {
	GenerateToken(string) (string, error)
	SecretKey() (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

type jwtManager struct {
	secretKey []byte
	logger    *logrus.Logger
}

func NewJWT(secretKey []byte, logger *logrus.Logger) JWTManager {
	return jwtManager{
		secretKey: secretKey,
		logger:    logger,
	}
}

func (s jwtManager) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": jwt.TimeFunc().Add(time.Hour * 24).Unix(),
	}

	s.logger.Info("Generating claims...")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := []byte(s.secretKey)
	tokenString, err := token.SignedString(key)
	if err != nil {
		s.logger.Errorf("token string was not generated, %s", err.Error())
		return "", err
	}

	s.logger.Info("Token string generated.")
	return tokenString, nil
}

func (s jwtManager) SecretKey() (string, error) {
	strSecretKey := string(s.secretKey)
	if strSecretKey == "" {
		msg := "Secret Key is empty"
		s.logger.Info(msg)
		return "", errors.New(msg)
	}

	return string(strSecretKey), nil
}

func (s jwtManager) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}
