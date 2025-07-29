package models

type AuthUser struct {
	ID       int     `db:"id" json:"id"`
	PersonID *int    `db:"person_id" json:"person_id"`
	Email    string  `db:"email" json:"email"`
	Password string  `db:"password_hash" json:"-"`
	Role     string  `db:"role" json:"role"`
	Person   *Person `json:"person,omitempty"`
}
