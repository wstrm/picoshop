package auth

import (
	"net/http"
	"time"

	"github.com/willeponken/picoshop/session"
)

const (
	cookieName = "auth"
)

type Authorizer struct {
	manager *session.Manager
}

func NewAuthorizer(maxLifeTime time.Duration) *Authorizer {
	return &Authorizer{
		manager: session.NewManager(cookieName, maxLifeTime),
	}
}

func (authorizer *Authorizer) Wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session := authorizer.manager.Start(writer, request)

		id := session.Get("id")
		if id == nil {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (authorizer *Authorizer) Login(id int64, writer http.ReponseWriter, request *http.Request) {
	authorizer.manager.Start(writer, request).Set("id", id)
}

func (authorizer *Authorizer) Logout(writer http.ReponseWriter, request *http.Request) {
	authorizer.manager.DestroySession(writer, request)
}
