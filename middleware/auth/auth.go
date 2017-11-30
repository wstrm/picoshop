package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/willeponken/picoshop/session"
)

const (
	cookieName = "auth"
)

func failedToAuthenticateSessionError() error {
	return errors.New("Failed to authenticate session")
}

func getEmail(request *http.Request) (string, error) {
	s, err := session.Get(request, cookieName)
	if err != nil {
		return "", err
	}

	if email, ok := s.Values["email"].(string); ok {
		return email, nil
	}

	return "", nil
}

func injectEmail(writer http.ResponseWriter, request *http.Request, next http.Handler, protected bool) {
	email, err := getEmail(request)
	if err != nil {
		log.Println(err)
		http.Error(writer, failedToAuthenticateSessionError().Error(),
			http.StatusInternalServerError)
		return
	}

	if email != "" {
		ctx := context.WithValue(request.Context(), "email", email)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	} else if protected {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	} else {
		next.ServeHTTP(writer, request)
	}
}

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		injectEmail(writer, request, next, true)
	})
}

func Intercept(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		injectEmail(writer, request, next, false)
	})
}

func Login(email string, writer http.ResponseWriter, request *http.Request) error {
	s, err := session.Get(request, cookieName)
	if err != nil {
		return err
	}

	s.Values["email"] = email
	session.Save(request, writer, s)

	return nil
}

func Logout(writer http.ResponseWriter, request *http.Request) error {
	s, err := session.Get(request, cookieName)
	if err != nil {
		return err
	}

	s.Options.MaxAge = -1 // Invalidate session
	session.Save(request, writer, s)

	return nil
}

func IsLoggedIn(request *http.Request) bool {
	s, _ := session.Get(request, cookieName)
	if _, ok := s.Values["email"]; ok {
		return true
	}

	return false
}
