package services

import (
	"fmt"
	"jusan_demo/pkg/db"
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

func (s *PersonService) UpsertPerson(req UpsertPersonRequest, authUserID int) error {
	var personID int
	stateCode := "P1"

	err := db.DB.QueryRow(`
		call upsert_person(
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		`, req.LastName, req.FirstName, req.BirthDate, req.PhoneNumber,
		req.UIN, stateCode, req.BranchID, req.PersonID, req.RedactorID,
	).Scan(&personID)
	if err != nil {
		fmt.Println("Ошибка при вызове upsert_person:", err)
		return err
	}

	fmt.Println("Получен personID из процедуры:", personID)

	_, err = db.DB.Exec(`update auth_user set person_id = $1 where id = $2`, personID, authUserID)
	if err != nil {
		fmt.Println("Ошибка при обновлении auth_user:", err)
		return err
	}

	//Создаём карту
	_, err = db.DB.Exec(`
		call upsert_card(
			null, 
			CURRENT_DATE + INTERVAL '3 years',
			'C1-VIR',
			'C1',
			$1, -- person_id
			null,
			CURRENT_DATE,
			true,
			true,
			null,
			$2 -- redactor_id
		)
	`, personID, req.RedactorID)
	if err != nil {
		fmt.Println("Ошибка при создании карты:", err)
	}

	return err
}
