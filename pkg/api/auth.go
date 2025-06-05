package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go_final_project/pkg/config"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		var storedPasword = config.Password

		var pass Password

		err := json.NewDecoder(r.Body).Decode(&pass)
		if err != nil {
			log.Println("JSON deserialization error", err)
			http.Error(w, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}

		if pass.Password == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "need to enter a password"})
			return
		}

		if pass.Password == storedPasword {
			token, err := GetToken(pass.Password)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"token": token})
		}

		if pass.Password != storedPasword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "incorrect password"})
			return
		}
	}
}

func GetToken(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	hashString := hex.EncodeToString(hash[:])

	claims := &Claims{
		PasswordHash: hashString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("failed to sign jwt:", err)
		return "", err
	}

	log.Println("Result token:", signedToken)

	return signedToken, nil
}

func ValidateToken(tokenString string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}

	if ok := token.Valid; ok {
		return true, nil
	}

	return false, nil
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var storedPasword = config.Password

		fmt.Println("Password:", storedPasword)

		if len(storedPasword) > 0 {
			var cookieToken string

			cookie, err := r.Cookie("token")
			if err == nil {
				cookieToken = cookie.Value
			}

			fmt.Println("JWT:", cookieToken)

			var isValid bool

			isValid, err = ValidateToken(cookieToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fmt.Println("validation successful")

			if !isValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	}
}
