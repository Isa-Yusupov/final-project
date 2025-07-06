package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

type ResponseData struct {
	Password string `json:"password"`
}

var secretKey = []byte("secret_key")

func generateToken(hash string) (string, error) {
	claims := jwt.MapClaims{
		"hash": hash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func validateToken(tokenString string, expectedHash string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["hash"].(string)
	return ok && hash == expectedHash
}

func passwordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"Некорректный запрос"}`, http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || body.Password != pass {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	hash := passwordHash(pass)
	token, err := generateToken(hash)
	if err != nil {
		http.Error(w, `{"error":"Ошибка сервера"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		if !validateToken(cookie.Value, passwordHash(pass)) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
