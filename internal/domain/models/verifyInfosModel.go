package models

type VerifyInfos struct {
	UserID            uint64
	verificationToken string
	IsVerified        bool
}
