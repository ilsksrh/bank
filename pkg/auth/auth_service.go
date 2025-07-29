package auth

import (
	"errors"
	"jusan_demo/pkg/config"
	"jusan_demo/pkg/db"
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
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя с email, паролем и ролью
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} map[string]string
// @Failure 409 {string} string "Email уже используется"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /register [post]
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


// Login godoc
// @Summary Вход пользователя
// @Description Авторизация по email и паролю, выдаёт access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Email и пароль"
// @Success 200 {object} TokenResponse
// @Failure 401 {string} string "Неверный пароль или пользователь"
// @Router /login [post]
// In pkg/auth/auth_service.go
func (s *AuthService) Login(email, password string) (string, string, error) {
    var (
        userID       int
        personID     int
        passwordHash string
        role         string
    )

    err := db.DB.QueryRow(`
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

// In pkg/auth/auth_service.go
func (s *AuthService) generateToken(userID, personID int, role, secret string, ttl time.Duration) (string, error) {
    // You'll need to get the email from the database when generating the token
    var email string
    err := db.DB.QueryRow("SELECT email FROM auth_user WHERE id = $1", userID).Scan(&email)
    if err != nil {
        return "", err
    }

    claims := TokenClaims{
        UserID:   userID,
        PersonID: personID,
        Role:     role,
        Email:    email, // Add the email to claims
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
