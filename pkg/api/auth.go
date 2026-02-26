package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

type signinRequest struct {
	Password string `json:"password"`
}

type jwtPayload struct {
	PH  string `json:"ph"`
	Exp int64  `json:"exp"`
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func signToken(password string) (string, error) {
	headerRaw, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	payloadRaw, err := json.Marshal(jwtPayload{
		PH:  passwordHash(password),
		Exp: time.Now().Add(8 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding
	header := enc.EncodeToString(headerRaw)
	payload := enc.EncodeToString(payloadRaw)
	unsigned := header + "." + payload

	mac := hmac.New(sha256.New, []byte(password))
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return "", err
	}
	signature := enc.EncodeToString(mac.Sum(nil))

	return unsigned + "." + signature, nil
}

func verifyToken(token, password string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	unsigned := parts[0] + "." + parts[1]
	enc := base64.RawURLEncoding

	mac := hmac.New(sha256.New, []byte(password))
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return false
	}
	wantSig := enc.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(wantSig)) {
		return false
	}

	payloadRaw, err := enc.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var payload jwtPayload
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return false
	}
	if payload.PH != passwordHash(password) {
		return false
	}
	if time.Now().Unix() >= payload.Exp {
		return false
	}
	return true
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil || !verifyToken(cookie.Value, pass) {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error())
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeError(w, errors.New("authentication is disabled").Error())
		return
	}
	if req.Password != pass {
		writeError(w, "invalid password")
		return
	}

	token, err := signToken(pass)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
