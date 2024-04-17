package model

import "github.com/google/uuid"

//TokenPair used for returning pairs of id refresh Tokens

//type TokenPair struct {
//	IDToken      string `json:"idToken"`
//	RefreshToken string `json:"refreshToken"`
//}

type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refreshToken"`
}

type IDToken struct {
	SS string `json:"idToken"`
}

type TokenPair struct {
	IDToken
	RefreshToken
}
