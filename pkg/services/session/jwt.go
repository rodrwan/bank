package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/twinj/uuid"
)

// Service Auth service to handle jwt serialization.
type Service struct {
	accessSecret  []byte
	refreshSecret []byte

	store *redis.Client
}

// GetAuthData extract data from token.
func (svc *Service) GetAuthData(ctx context.Context, token string) (*Data, error) {
	claims := &Claims{}
	data, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return svc.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	val, ok := data.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return val.Data, nil
}

// CreateAuthData Creates a new Auth.
func (svc *Service) CreateAuthData(ctx context.Context, referenceID string, payload []byte) (*Auth, error) {
	tokens, err := CreateTokensWithSecrets(referenceID, string(payload), svc.accessSecret, svc.refreshSecret)
	if err != nil {
		return nil, err
	}

	// insert new token into db
	if err := svc.CreateAuth(referenceID, tokens); err != nil {
		return nil, err
	}

	return &Auth{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// CreateAuth insert token into store
func (svc *Service) CreateAuth(referenceID string, ad *Data) error {
	at := time.Unix(ad.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(ad.RtExpires, 0)
	now := time.Now()

	atCreated, err := svc.store.Set(ad.AccessUUID, referenceID, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := svc.store.Set(ad.RefreshUUID, referenceID, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}

	return nil
}

// RefreshAuthData Update and existing Auth.
func (svc *Service) RefreshAuthData(ctx context.Context, auth *Auth) (*Auth, error) {
	// TODO: get refresh token from context
	// verify refresh token
	// token is valid?
	// get refresh token uuid
	// delete refresh token by uuid
	// create new access_token and refresh token
	return nil, errors.New("not implemented")
}

// BlockAuthData Block an existing Auth, this is intended to prevent future hack or stolen token
// this will be store jwt metadata in a redis database.
func (svc *Service) BlockAuthData(ctx context.Context, auth *Auth) error {
	// get access_token from context
	// remove tokens from store
	return nil
}

// Data represets a jwt token metadata
type Data struct {
	Data         string `json:"data,omitempty"`
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
	AccessUUID   string `json:"access_uuid,omitempty"`
	RefreshUUID  string `json:"refresh_uuid,omitempty"`
	AtExpires    int64  `json:"-"`
	RtExpires    int64  `json:"-"`
	ReferenceID  string
}

// Auth represents auth data.
type Auth struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
}

// Claims join Data and StandardClaims from jwt.
type Claims struct {
	jwt.StandardClaims
	Data        *Data  `json:"data"`
	ReferenceID string `json:"reference_id"`
}

// CreateTokensWithSecrets creates a new access_token and refresh_token with the given auth data and secret.
func CreateTokensWithSecrets(referenceID, data string, aSecret, rSecret []byte) (*Data, error) {
	iat := time.Now()
	aExp := time.Now().Add(time.Minute * 15)
	rExp := time.Now().Add(time.Hour * 24 * 7)

	jwtData := &Data{
		AccessUUID:  uuid.NewV4().String(),
		RefreshUUID: uuid.NewV4().String(),
		AtExpires:   aExp.Unix(),
		RtExpires:   rExp.Unix(),
	}

	aStdClms := jwt.StandardClaims{
		Id:        uuid.NewV4().String(),
		IssuedAt:  iat.Unix(),
		ExpiresAt: aExp.Unix(),
	}
	atClaims := &Claims{
		StandardClaims: aStdClms,
		Data: &Data{
			AccessUUID: jwtData.AccessUUID,
			Data:       data,
		},
		ReferenceID: referenceID,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(aSecret)
	if err != nil {
		return nil, err
	}

	jwtData.AccessToken = accessToken

	rStdClms := jwt.StandardClaims{
		Id:        uuid.NewV4().String(),
		IssuedAt:  iat.Unix(),
		ExpiresAt: rExp.Unix(),
	}
	rtClaims := &Claims{
		StandardClaims: rStdClms,
		Data: &Data{
			RefreshUUID: jwtData.RefreshUUID,
		},
		ReferenceID: referenceID,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(rSecret)
	if err != nil {
		return nil, err
	}

	jwtData.RefreshToken = refreshToken

	return jwtData, nil
}
