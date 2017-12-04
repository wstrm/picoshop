// Package session wraps a gorilla/sessions cookie store. In the future this
// could potentially be replaced by a in-house solution (per the Picoshop philosophy).
package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte(generateSecret()))
}

func generateSecret() string {
	b := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

// Get returns a session using its name.
func Get(request *http.Request, name string) (*sessions.Session, error) {
	return store.Get(request, name)
}

// New creates a new session by a provided name.
func New(request *http.Request, name string) (*sessions.Session, error) {
	return store.New(request, name)
}

// Save writes a new session.
func Save(request *http.Request, writer http.ResponseWriter, session *sessions.Session) error {
	return store.Save(request, writer, session)
}

// MaxAge sets the session max age.
func MaxAge(age int) {
	store.MaxAge(age)
}
