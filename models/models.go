package models

import "gorm.io/gorm"

var DB *gorm.DB

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}
type Role struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
type Permission struct {
	gorm.Model
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}
