package utils

import (
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type jwtCustomClaims struct {
	repository.Profile
	jwt.RegisteredClaims
}

func GenerateToken(dataUser repository.Profile) (t, expiresAtStr string, err error) {
	expiresAt := time.Now().Add(time.Hour * 72)
	expiresAtStr = expiresAt.Format("2006-01-02 15:04:04")

	prvKey, err := ioutil.ReadFile("secret_cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		Profile: repository.Profile{
			UserId:   dataUser.UserId,
			FullName: dataUser.FullName,
			Phone:    dataUser.Phone,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Create token with claims
	t, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	return
}

func ValidateToken(tokenReq string) (claims interface{}, err error) {
	tokenString := strings.Replace(tokenReq, "Bearer ", "", -1)

	pubKey, err := ioutil.ReadFile("secret_cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return
	}

	tok, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		err = errors.New("invalid token")
		return
	}

	return
}
