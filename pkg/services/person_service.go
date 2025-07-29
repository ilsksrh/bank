package services

import (
	"fmt"
	"jusan_demo/pkg/db"
	"time"
)

type PersonService struct{}

func NewPersonService() *PersonService {
	return &PersonService{}
}

type UpsertPersonRequest struct {
	LastName    string  `json:"last_name"`
	FirstName   string  `json:"first_name"`
	BirthDate   string  `json:"date_of_birth"`
	PhoneNumber string  `json:"phone_number"`
	UIN         string  `json:"uin"`
	StateID     *string `json:"state_id"`
	BranchID    *int    `json:"branch_id"`
	PersonID    *int    `json:"person_id"`
	RedactorID  *int    `json:"redactor_id"`
}

// @Summary Создание/обновление профиля с выпуском карты
// @Description После создания профиля создаётся аккаунт и виртуальная карта
// @Tags person
// @Accept json
// @Produce json
// @Param request body services.UpsertPersonRequest true "Данные пользователя"
// @Success 200 {string} string "Зарегистрирован"
// @Failure 400 {string} string "Ошибка запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/person/upsert [post]
func (s *PersonService) UpsertPerson(req UpsertPersonRequest, authUserID int, isAdmin bool) error {
    var personID int
    stateCode := "P1"

    // Parse birth date
    birthDate, err := time.Parse("2006.01.02", req.BirthDate)
    if err != nil {
        return fmt.Errorf("invalid date format for birth date, expected YYYY.MM.DD: %w", err)
    }

    // Determine redactor_id based on operation type
    var redactorID *int
    if req.PersonID != nil { // This is an update operation
        redactorID = &authUserID // Always use authUserID as redactor for updates
    } 
    // For new profiles, redactorID remains nil

    // 1. First create/update the person record
    err = db.DB.QueryRow(
        `SELECT * FROM fn_upsert_person(
            $1, $2, $3::date, $4, $5, $6, $7, $8, $9
        )`,
        req.LastName, 
        req.FirstName, 
        birthDate.Format("2006-01-02"), // Convert to SQL date format
        req.PhoneNumber,
        req.UIN, 
        stateCode, 
        req.BranchID, 
        req.PersonID, 
        redactorID,
    ).Scan(&personID)

    if err != nil {
        return fmt.Errorf("ошибка создания/обновления профиля: %w", err)
    }

    // 2. Link the person to auth_user (only for new profiles)
    if req.PersonID == nil {
        _, err = db.DB.Exec(
            `UPDATE auth_user SET person_id = $1 WHERE id = $2`,
            personID, authUserID,
        )
        if err != nil {
            return fmt.Errorf("ошибка привязки профиля к пользователю: %w", err)
        }

        // 3. Create account and card (only for new profiles)
        // Create account
        var accountResponse struct {
            AccountID     int    `db:"account_id"`
            AccountNumber string `db:"account_number"`
        }
        
        err = db.DB.QueryRowx(
            `SELECT * FROM fn_upsert_account(
                $1, NULL, $2, $3, 0, $4
            )`,
            personID, "A1-D", "A1", &authUserID,
        ).StructScan(&accountResponse)
        
        if err != nil {
            return fmt.Errorf("ошибка создания счёта: %w", err)
        }

        // Create card
        issueDate := time.Now()
        expirationDate := issueDate.AddDate(3, 0, 0) // 3 years from now
        
        var cardResponse struct {
            CardID     int    `db:"card_id"`
            CardNumber string `db:"card_number"`
        }
        
        err = db.DB.QueryRowx(
            `SELECT * FROM fn_upsert_card(
                $1::date,  -- expiration_date
                $2,        -- type_code
                $3,        -- state_code
                $4,        -- person_id
                $5,        -- account_id
                NULL,      -- card_number (let function generate)
                $6::date,  -- issue_date
                true,      -- is_default
                true,      -- is_virtual
                NULL,      -- card_id
                $7         -- redactor_id
            )`,
            expirationDate,
            "C1-VIR", 
            "C1", 
            personID,
            accountResponse.AccountID, 
            issueDate,
            &authUserID,
        ).StructScan(&cardResponse)
        
        if err != nil {
            return fmt.Errorf("ошибка создания карты: %w", err)
        }
    }

    return nil
}

