package session

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var secretKey = "HalloweenIsComingSoonAndIAmGladAboutIt"
var SessionKey = "user"
var (
	ErrBadSignMethod = errors.New("bad sign method")
	ErrBadToken      = errors.New("bad token")
	ErrNoPayload     = errors.New("no payload")
)

func Create(id, login string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": login,
			"id":       id,
		},
	})
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Check(inToken string) (user map[string]interface{}, err error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, ErrBadSignMethod
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, ErrBadToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrNoPayload
	}
	user = payload["user"].(map[string]interface{})
	return user, nil
}

func SessionFromContext(ctx context.Context) (string, string) {
	user := ctx.Value(SessionKey).(map[string]interface{})
	id := user["id"].(string)
	login := user["username"].(string)
	return id, login
}
