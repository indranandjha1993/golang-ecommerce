package jwt

import (
	"fmt"
	"time"

	"golang-ecommerce/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

type JWT struct {
	secret string
	db     *gorm.DB
}

func New(secret string, db *gorm.DB) *JWT {
	return &JWT{secret: secret, db: db}
}

func (j *JWT) GenerateToken(userID uint, expiresIn time.Duration) (*Token, error) {
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
		},
		UserID: userID,
	})

	// Sign and get the complete encoded token as a string
	accessToken, _ := token.SignedString([]byte(j.secret))

	// Create a new refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// Sign and get
	refreshTokenString, _ := refreshToken.SignedString([]byte(j.secret))

	// Save the tokens to the UserToken model
	userToken := &models.UserToken{
		UserID:    userID,
		Token:     accessToken,
		ExpiresAt: time.Now().Add(expiresIn).Unix(),
		Revoked:   false,
	}
	if err := j.db.Create(userToken).Error; err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    time.Now().Add(expiresIn).Unix(),
	}, nil
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	// Get the Claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token is revoked
	var userToken models.UserToken
	if err := j.db.Where("token = ?", tokenString).First(&userToken, "user_id = ?", claims.UserID).Error; err != nil {
		return nil, err
	}
	if userToken.Revoked {
		return nil, fmt.Errorf("token has been revoked")
	}

	return claims, nil
}
