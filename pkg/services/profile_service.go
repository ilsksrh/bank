package services

import (
	"jusan_demo/pkg/db"
	"jusan_demo/pkg/models"
	"fmt"
)

type ProfileService struct{}

func NewProfileService() *ProfileService {
	return &ProfileService{}
}
// services/profile_service.go
func (s *ProfileService) GetProfile(userID int) (*models.AuthUser, error) {
    var user models.AuthUser
    err := db.DB.Get(&user, `
        SELECT id, email, role, person_id 
        FROM auth_user 
        WHERE id = $1`, userID)
    if err != nil || user.PersonID == nil {
        return nil, fmt.Errorf("user not found or not linked to profile")
    }

    var person models.Person
    err = db.DB.Get(&person, `
        SELECT id, first_name, last_name, date_of_birth, uin, phone_number 
        FROM person 
        WHERE id = $1`, *user.PersonID)
    if err != nil {
        return nil, fmt.Errorf("failed to get profile: %w", err)
    }

    var accounts []models.Account
    err = db.DB.Select(&accounts, `
        SELECT id, account_number, account_balance 
        FROM account 
        WHERE person_id = $1`, person.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get accounts: %w", err)
    }

    for i := range accounts {
        var cards []models.Card
        err = db.DB.Select(&cards, `
            SELECT id, card_number, is_virtual, is_default, issue_date, expiration_date 
            FROM card 
            WHERE account_id = $1`, accounts[i].ID)
        if err != nil {
            return nil, fmt.Errorf("failed to get cards: %w", err)
        }
        accounts[i].Cards = cards
    }

    person.Accounts = accounts
    user.Person = &person

    return &user, nil
}