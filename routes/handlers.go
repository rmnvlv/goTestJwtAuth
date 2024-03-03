package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Session struct {
	GUID         string
	RefreshToken string
}

func GetTokens(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "aplication/json")

	var user Session
	json.NewDecoder(request.Body).Decode(&user)
	fmt.Printf("User: %v\n", user)

	userDB := SelectUserFromDbByGuid(user.GUID)

	if userDB.GUID != "" {
		tokens, err := createTokens(user.GUID)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			fmt.Printf("Eroor with creating token: %v", err)
		}

		err = InsertInDB(user.GUID, tokens["refresh-token"])
		if err != nil {
			log.Fatal(err)
		}

		writer.Header().Set("Authorization", fmt.Sprintf("Bearer %v", tokens["access-token"]))
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, tokens["access-token"])
		fmt.Printf("Access-token: %v\n", tokens["access-token"])
		fmt.Printf("Refresh-token: %v\n", tokens["refresh-token"])
		return
	} else {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprint(writer, "Invalid credentials", userDB)
	}
}

func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "aplication/json")

	var user Session
	json.NewDecoder(request.Body).Decode(&user)
	fmt.Printf("User: %v\n", user)

	//Проверить являетсяли токен валидным
	if verifyRefreshToken(user.RefreshToken) != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(writer, "Invalid token")
		return
	}

	userDB := SelectSessionFromDbByGuid(user.GUID)

	if userDB.RefreshToken != "" && user.RefreshToken == userDB.RefreshToken {
		tokens, err := createTokens(user.GUID)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			fmt.Printf("Eroor with creating token: %v", err)
		}

		err = UpdateSession(userDB.GUID, tokens["refresh-token"])
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			fmt.Printf("Eroor with updating session: %v", err)
		}

		writer.Header().Set("Authorization", fmt.Sprintf("Bearer %v", tokens["access-token"]))
		writer.WriteHeader(http.StatusOK)

		fmt.Fprint(writer, tokens["refresh-token"])
		fmt.Printf("Token: %v\n", tokens)
		//change refresh token in db
		return
	} else {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprint(writer, "Invalid credentials")
	}
}

func Protected(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(writer, "Missing authorization header")
		return
	}

	tokenString = tokenString[len("Bearer "):]
	if verifyAccessToken(tokenString) != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(writer, "Invalid token")
		return
	}

	fmt.Fprint(writer, "Welcome to the the protected area")
}
