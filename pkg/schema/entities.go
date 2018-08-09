package schema

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type User struct {
	Id    string `json:"id" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email string `json:"email" gorm:"type:varchar(100);unique_index"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Id", uuid.NewV4().String())
	return nil
}
