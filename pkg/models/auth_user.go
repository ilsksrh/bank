package models

type AuthUser struct {
    ID       int    `db:"id" json:"-"`
    PersonID *int   `db:"person_id"` 
    Email    string `db:"email"`
    Password string `db:"password_hash" json:"-"`
    Role     string `db:"role"`
	
}