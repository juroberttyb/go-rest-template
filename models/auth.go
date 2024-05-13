package models

import jwt "github.com/golang-jwt/jwt/v5"

type SourceAppID string

const (
// OTHER_APP_1 SourceAppID = "f35edd48-78c1-413f-a171-a775cb6defe3"
// OTHER_APP_2 SourceAppID = "..."
)

type Claims struct {
	jwt.RegisteredClaims
	Scope UserType `json:"scope"`
}
