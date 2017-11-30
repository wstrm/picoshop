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

func Get(request *http.Request, name string) (*sessions.Session, error) {
	return store.Get(request, name)
}

func New(request *http.Request, name string) (*sessions.Session, error) {
	return store.New(request, name)
}

func Save(request *http.Request, writer http.ResponseWriter, session *sessions.Session) error {
	return store.Save(request, writer, session)
}

func MaxAge(age int) {
	store.MaxAge(age)
}
