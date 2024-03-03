package routes

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKeyAccess = []byte("124312dsrfw34452dwqwe12")
var secretKeyRefresh = []byte("1010")

func createTokens(guid string) (map[string]string, error) {
	//gen token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"GUID": guid,
			"exp":  time.Now().Add(time.Minute * 5).Unix(),
		})

	t, err := token.SignedString(secretKeyAccess)
	if err != nil {
		return map[string]string{}, err
	}

	//gen refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 36).Unix(),
		})

	rt, err := refreshToken.SignedString(secretKeyRefresh)
	if err != nil {
		return map[string]string{}, err
	}

	tokens := map[string]string{
		"access-token":  t,
		"refresh-token": rt,
	}

	return tokens, nil
}

func verifyAccessToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKeyAccess, nil
	})
	// fmt.Printf("Token on veryfy: %v\n: ", token)

	if err != nil {
		return err
	}

	if !t.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func verifyRefreshToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKeyRefresh, nil
	})
	// fmt.Printf("Token on veryfy: %v\n: ", t)

	if err != nil {
		return err
	}

	if !t.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
