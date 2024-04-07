package model

//TokenPair used for returning pairs of id refresh Tokens

type TokenPair struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}
