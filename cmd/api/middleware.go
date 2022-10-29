package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

//Simply wrapping the signal with CORS allowence
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			//could set as anon user;
		}
		headerparts := strings.Split(authHeader, " ")

		if len(headerparts) != 2 {
			app.errorJson(w, errors.New("invalid auth header"))
			return
		}

		if headerparts[0] != "Bearer" {
			app.errorJson(w, errors.New("unauthorized-no Bearer"))
			return
		}

		token := headerparts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJson(w, errors.New("unauthorized-Hmacc check failed!"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJson(w, errors.New("unauthorized-token expired!"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.errorJson(w, errors.New("unauthorized-invalid audience!"), http.StatusForbidden)
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.errorJson(w, errors.New("unauthorized-invalid issuer!"), http.StatusForbidden)
			return
		}
		UserID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJson(w, errors.New("unauthorized!"), http.StatusForbidden)
			return
		}
		log.Println("Valid user: ", UserID)
		next.ServeHTTP(w, r)
	})
}
