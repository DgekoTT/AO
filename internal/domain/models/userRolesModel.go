package models

type UserRole string

type UserRoles struct {
	UserID uint64
	Role   UserRole
}
