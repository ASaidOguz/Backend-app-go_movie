package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	ID:       10,
	Email:    "ahmetogus123@gmail.com",
	Password: "$2a$12$zCJN9bDE0Pog17vErHrFQeZjKXZtog6P/dZPs.GQekZHvO5WMl4Qq",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {

	var cred Credentials
	cred.Password = "$2a$12$zCJN9bDE0Pog17vErHrFQeZjKXZtog6P/dZPs.GQekZHvO5WMl4Qq"

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		app.errorJson(w, errors.New("unauthorized!!!"))
		return
	}

	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(cred.Password))
	if err != nil {
		fmt.Println(cred.Username, cred.Password)
		app.errorJson(w, errors.New("unauthorized!!!"))
		return
	}
	var claims jwt.Claims

	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {

		app.errorJson(w, errors.New("error signing in!!!"))
		return
	}
	app.WriteJSON(w, http.StatusOK, jwtBytes, "response")
	fmt.Println(jwtBytes)
}
