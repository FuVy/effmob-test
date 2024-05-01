package storage

import "github.com/google/uuid"

// swagger:model CarObject
type Car struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;autoIncrement:false" json:"id"`
	RegNum          string    `json:"reg_num" example:"X123XX150" gorm:"uniqueIndex"`
	Mark            string    `json:"mark" example:"Lada"`
	Model           string    `json:"model" example:"Vesta"`
	Year            *int      `json:"year" example:"2002"`
	OwnerName       string    `json:"owner_name"`
	OwnerSurname    string    `json:"owner_surname"`
	OwnerPatronymic *string   `json:"owner_patronymic,omitempty"`
} // @Name CarObject
