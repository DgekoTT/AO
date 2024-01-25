package models

import "github.com/docker/distribution/uuid"

type Users struct {
	ID                   uuid.UUID
	Status               string
	RegistrationProvider string
	Roles                []*UserRoles
	VerifyInfos          VerifyInfos
	AccessToken          string
	RefreshToken         string
	OauthProviders       []*OAuthProviders
	Email                string
	HashedPassword       []byte
}
