package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/session"
)

const (
	cookieName = "auth"
	userKey    = "user"
)

func init() {
	gob.Register(model.User{})
}

func failedToAuthenticateSessionError() error {
	return errors.New("Failed to authenticate session")
}

func getUser(request *http.Request) (model.User, error) {
	s, err := session.Get(request, cookieName)
	if err != nil {
		return model.User{}, err
	}

	if user, ok := s.Values[userKey].(model.User); ok {
		return user, nil
	}

	return model.User{}, nil
}

func injectUser(writer http.ResponseWriter, request *http.Request, next http.Handler, protected bool) {
	user, err := getUser(request)
	if err != nil {
		http.Error(writer, failedToAuthenticateSessionError().Error(),
			http.StatusInternalServerError)
		return
	}

	if !user.IsValid() {
		ctx := context.WithValue(request.Context(), userKey, user)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	} else if protected {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	} else {
		next.ServeHTTP(writer, request)
	}
}

func Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		injectUser(writer, request, next, true)
	})
}

func Intercept(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		injectUser(writer, request, next, false)
	})
}

func Login(user model.User, writer http.ResponseWriter, request *http.Request) error {
	s, err := session.Get(request, cookieName)
	if err != nil {
		return err
	}

	s.Values[userKey] = user
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
	if _, ok := s.Values[userKey]; ok {
		return true
	}

	return false
}
