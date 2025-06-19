package auth

import (
	"errors"
	"jusan_demo/pkg/config"
	"jusan_demo/pkg/db"
	"jusan_demo/pkg/models"
	"time"
	"github.com/lib/pq"
	"log"


	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	UserID   int    `json:"user_id"`
	PersonID int    `json:"person_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(email, password, role string, personID *int) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Ошибка хеширования пароля:", err)
		return err
	}

	_, err = db.DB.Exec(`
		INSERT INTO auth_user (email, password_hash, role)
		VALUES ($1, $2, $3)
	`, email, string(hashedPassword), role)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Println("Ошибка: email уже используется:", email)
			return errors.New("email уже используется")
		}
		log.Println("Ошибка записи в БД:", err)
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	var (
		userID       int
		personID     int
		passwordHash string
		role         string
	)

	err := db.DB.QueryRowx(`
		SELECT id, password_hash, role
		FROM auth_user
		WHERE email = $1
	`, email).Scan(&userID, &passwordHash, &role)
	if err != nil {
		return "", "", errors.New("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", "", errors.New("неверный пароль")
	}

	access, err := s.generateToken(userID, personID, role, config.AppConfig.AccessSecret, config.AppConfig.AccessTTL)
	if err != nil {
		return "", "", err
	}

	refresh, err := s.generateToken(userID, personID, role, config.AppConfig.RefreshSecret, config.AppConfig.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) generateToken(userID, personID int, role, secret string, ttl time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:   userID,
		PersonID: personID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) VerifyAccessToken(tokenStr string) (*TokenClaims, error) {
	return s.verifyToken(tokenStr, config.AppConfig.AccessSecret)
}

func (s *AuthService) VerifyRefreshToken(tokenStr string) (*TokenClaims, error) {
	return s.verifyToken(tokenStr, config.AppConfig.RefreshSecret)
}

func (s *AuthService) verifyToken(tokenStr, key string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) GetProfile(userID int) (*models.AuthUser, error) {
	var user models.AuthUser
	err := db.DB.Get(&user, `SELECT id, email, role, person_id FROM auth_user WHERE id = $1`, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}
	return &user, nil
}
