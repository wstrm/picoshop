package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"reflect"

	"github.com/willeponken/picoshop/session"
)

type userKey string

type User interface {
	IsValid() bool
}

type Manager struct {
	cookieName string
	// TODO(willeponken): benchmark vs. something like map[userKey]bool
	supportedUsers []userKey // short array - should be faster than hash map (?)
}

type Policy struct {
	authorization map[userKey]bool
	protected     bool
}

func NewPolicy() Policy {
	return Policy{}
}

func (policy Policy) SetUser(user User, authorize bool) {
	policy.protected = true

	if policy.authorization == nil {
		policy.authorization = make(map[userKey]bool)
	}
	policy.authorization[getUserKey(user)] = authorize
}

func (policy Policy) SetProtected(protected bool) {
	if len(policy.authorization) > 0 && !protected {
		panic("policy cannot be unprotected and contain authorization map for users")
	}
	policy.protected = protected
}

func (policy Policy) isAllowed(key userKey) bool {
	if !policy.protected {
		return true // unprotected
	}

	if authorization, ok := policy.authorization[key]; ok {
		return authorization // protected, return if user is allowed
	}

	return false // protected, but no such user was found
}

func NewManager(cookieName string, userInterfaces ...User) *Manager {
	var supportedUsers []userKey

	for _, userInterface := range userInterfaces {
		gob.Register(userInterface)                                        // will panic if registered multiple managers with same user interface
		supportedUsers = append(supportedUsers, getUserKey(userInterface)) // hnnnnnggggggggg 😍
	}

	return &Manager{
		cookieName:     cookieName,
		supportedUsers: supportedUsers,
	}
}

func (manager *Manager) getUser(request *http.Request) (user User, err error) {
	s, err := session.Get(request, manager.cookieName)
	if err != nil {
		return
	}

	for supportedUser := range manager.supportedUsers {
		if u, ok := s.Values[supportedUser].(User); ok {
			// return first found user in session. Several users should never exist
			// and the Manager.Login method should make sure to empty old users before
			// inserting a new user.
			user = u
			return
		}
	}

	err = errors.New("cannot find any supported user inside session")
	return
}

func (manager *Manager) Middleware(next http.Handler, policy Policy) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := manager.getUser(request)

		// no user found, and unprotected handler
		if err != nil && !policy.protected {
			next.ServeHTTP(writer, request)
			return
		}

		key := getUserKey(user)

		// user found, add to context
		if user.IsValid() {
			ctx := context.WithValue(request.Context(), key, user)

			if policy.isAllowed(key) { // is allowed
				next.ServeHTTP(writer, request.WithContext(ctx))
			} else { // authorization declined
				http.Redirect(writer, request.WithContext(ctx), "/", http.StatusSeeOther)
			}
		}
	})
}

func (manager *Manager) Login(user User, writer http.ResponseWriter, request *http.Request) error {
	s, err := session.Get(request, manager.cookieName)
	if err != nil {
		return err
	}

	for supportedUser := range manager.supportedUsers {
		delete(s.Values, supportedUser) // in Go, delete(m, k) is no-op if key does not exist in map
	}

	s.Values[getUserKey(user)] = user
	session.Save(request, writer, s)

	return nil
}

func (manager *Manager) Logout(writer http.ResponseWriter, request *http.Request) error {
	s, err := session.Get(request, manager.cookieName)
	if err != nil {
		return err
	}

	s.Options.MaxAge = -1 // invalidate session - web browser will remove all entries
	session.Save(request, writer, s)

	return nil
}

func (manager *Manager) IsLoggedIn(request *http.Request) bool {
	s, _ := session.Get(request, manager.cookieName)

	// iterate over all keys in session to find 'any' registered user (should at most
	// be one due to Manager.Login)
	for supportedUser := range manager.supportedUsers {
		if _, ok := s.Values[supportedUser]; ok {
			return true // user key found
		}
	}

	return false
}

func failedToAuthenticateSessionError() error {
	return errors.New("Failed to authenticate session")
}

func getUserKey(user User) userKey {
	return userKey(reflect.TypeOf(user).Name())
}
