package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"reflect"

	"github.com/wstrm/picoshop/session"
)

type User interface {
	IsValid() bool
}

type Manager struct {
	cookieName string
	// TODO(wstrm): benchmark vs. something like map[userKey]bool
	supportedUsers []string // short array - should be faster than hash map (?)
}

type Policy struct {
	authorization map[string]bool
	protected     bool
}

func NewPolicy() Policy {
	return Policy{}
}

func (policy *Policy) SetUser(user User, authorize bool) {
	policy.protected = true

	if policy.authorization == nil {
		policy.authorization = make(map[string]bool)
	}
	policy.authorization[getUserKey(user)] = authorize
}

func (policy *Policy) SetProtected(protected bool) {
	if len(policy.authorization) > 0 && !protected {
		panic("policy cannot be unprotected and contain authorization map for users")
	}
	policy.protected = protected
}

func (policy *Policy) isAllowed(key string) bool {
	if !policy.protected {
		return true // unprotected
	}

	if authorization, ok := policy.authorization[key]; ok {
		return authorization // protected, return if user is allowed
	}

	return false // protected, but no such user was found
}

func NewManager(cookieName string, userInterfaces ...User) *Manager {
	var supportedUsers []string

	for _, userInterface := range userInterfaces {
		gob.Register(userInterface)                                        // will panic if registered multiple managers with same user interface
		supportedUsers = append(supportedUsers, getUserKey(userInterface)) // hnnnnnggggggggg 😍
	}

	return &Manager{
		cookieName:     cookieName,
		supportedUsers: supportedUsers,
	}
}

func (manager *Manager) GetUser(request *http.Request) (user User, err error) {
	s, err := session.Get(request, manager.cookieName)
	if err != nil {
		return
	}

	for _, supportedUser := range manager.supportedUsers {
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
		user, err := manager.GetUser(request)

		// no user found
		if err != nil || !user.IsValid() {
			if !policy.protected { // unprotected handler, allow
				next.ServeHTTP(writer, request)
				return
			} else { // protected, disallow
				http.Redirect(writer, request, "/", http.StatusSeeOther)
				return
			}
		}

		// user found, add to context
		key := getUserKey(user)
		ctx := context.WithValue(request.Context(), key, user)

		if policy.isAllowed(key) { // is allowed
			next.ServeHTTP(writer, request.WithContext(ctx))
		} else { // authorization declined
			http.Redirect(writer, request.WithContext(ctx), "/", http.StatusSeeOther)
		}
	})
}

func (manager *Manager) Login(user User, writer http.ResponseWriter, request *http.Request) error {
	s, _ := session.Get(request, manager.cookieName) // ignore errors, if it occurs the values will be overwritten anyhow

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

	session.Invalidate(request, writer, s)

	return nil
}

func (manager *Manager) IsLoggedIn(request *http.Request) bool {
	s, err := session.Get(request, manager.cookieName)
	if err != nil {
		return false
	}

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

func getUserKey(user User) string {
	return reflect.TypeOf(user).Name()
}
