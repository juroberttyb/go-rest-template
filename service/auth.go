package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"strings"
	"time"

	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/store"
	"github.com/A-pen-app/logging"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authSvc struct {
	c         store.Crypto
	verifyKey *rsa.PublicKey
}

var systemKeyID string
var app string

// NewAuth returns an implementation of service.Auth
func NewAuth(ctx context.Context, c store.Crypto) Auth {
	app = config.GetString("PROJECT_NAME")
	systemKeyID = config.GetString("SYSTEM_KEY_ID")
	base64Key, err := c.GetPublicKey(ctx, systemKeyID)
	if err != nil {
		panic(err)
	}
	key, err := base64.URLEncoding.DecodeString(base64Key)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		panic(err)
	}
	return &authSvc{
		c:         c,
		verifyKey: publicKey,
	}
}

type decryptedSignUpParams struct {
	AppID       string   `json:"app_id"`
	AppUserID   string   `json:"app_user_id"`
	Name        string   `json:"username"`
	Picture     *string  `json:"picture"`
	PushToken   *string  `json:"push_token"`
	Speciaities []string `json:"speciaities"`

	/* NOTE we don't include extra attributes because of limitation:
	 	crypto/rsa: message too long for RSA key size

	Gender     *models.Gender   `json:"gender"`
	Role       *models.UserRole `json:"role"`
	Facility   *string          `json:"facility"`
	Department *string          `json:"department"`
	Position   *string          `json:"position"`
	*/
}

// IssueToken returns a JWT for given user
func (a *authSvc) IssueToken(ctx context.Context, userID string, userType models.UserType, options ...IssueOption) (string, error) {
	opt := issueOption{
		ttl: 10 * 365 * 24 * time.Hour,
	}
	for _, f := range options {
		if err := f(&opt); err != nil {
			return "", err
		}
	}
	tokenID := uuid.New().String()
	claims := models.Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opt.ttl)),
			// TODO move string constant to config
			Issuer:   app,
			ID:       tokenID,
			Subject:  userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Audience: []string{userID},
		},
		userType,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	payload, err := t.SigningString()
	if err != nil {
		logging.Errorw(ctx, "create jwt payload failed", "err", err, "token", t.Raw)
		return "", err
	}
	sig, err := a.c.Sign(ctx, systemKeyID, payload)
	if err != nil {
		logging.Errorw(ctx, "sign jwt payload failed", "err", err, "payload", payload)
		return "", err
	}

	// remove trailing '='
	sig = strings.TrimRight(sig, "=")
	return strings.Join([]string{payload, sig}, "."), nil

}

// ValidateToken verifies if given token is issued by us
func (a *authSvc) ValidateToken(ctx context.Context, token string) (*models.Claims, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return a.verifyKey, nil
	}
	claims := models.Claims{}
	t, err := jwt.ParseWithClaims(token, &claims, keyFunc, jwt.WithIssuer(app))
	if err != nil {
		logging.Errorw(ctx, "parse jwt claims failed", "err", err, "token", token)
		return nil, err
	}
	if !t.Valid {
		logging.Errorw(ctx, "token invalid", "token", t.Raw)
		return nil, models.ErrorWrongParams
	}
	return &claims, nil
}
