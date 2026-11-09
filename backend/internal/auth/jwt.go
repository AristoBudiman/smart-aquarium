package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenDuration menentukan masa berlaku token JWT (24 jam sesuai kontrak).
const TokenDuration = 24 * time.Hour

// Claims mendefinisikan payload data pengguna yang disimpan di dalam JWT.
type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"sub"`
	Admin    bool   `json:"adm"`
	jwt.RegisteredClaims
}

// Service mengelola penandatanganan, verifikasi token JWT, dan hashing kata sandi.
type Service struct {
	secret []byte
}

// NewService membuat instance auth.Service baru dengan memverifikasi panjang kunci rahasia.
func NewService(secret string) (*Service, error) {
	if len(secret) < 32 {
		return nil, errors.New("auth: JWT_SECRET harus memiliki panjang minimal 32 karakter")
	}
	return &Service{secret: []byte(secret)}, nil
}

// Sign membuat dan menandatangani token JWT HS256 baru yang berlaku selama 24 jam.
func (s *Service) Sign(userID int64, username string, isAdmin ...bool) (string, error) {
	admin := len(isAdmin) > 0 && isAdmin[0]
	now := time.Now().UTC()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Admin:    admin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("auth: gagal menandatangani token: %w", err)
	}
	return signed, nil
}

// Verify memvalidasi token JWT HS256 dan mengekstrak klaim pengguna di dalamnya.
func (s *Service) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth: metode penandatanganan tidak valid: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth: token tidak valid: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("auth: klaim token tidak valid")
	}

	return claims, nil
}

// HashPassword menghasilkan hash bcrypt dari kata sandi teks biasa.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("auth: gagal melakukan hash kata sandi: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword memverifikasi apakah kata sandi teks biasa cocok dengan hash bcrypt.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
