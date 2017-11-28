package session

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/url"
	"sync"
	"time"
)

var globalSessions *Manager
var globalStorage *Storage

func init() {
	sessions = &Manager{
		storage: &Storage{},
	}
}

type Manager struct {
	cookieName   string
	sessionMutex sync.Mutex
	storage      *Storage
	maxLifeTime  int64
}

type Storage struct {
	mutex    sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

type Session struct {
	id         string
	lastAccess time.Time
	value      map[interface{}]interface{}
}

func (storage *Storage) Init(sessionId string) (Session, error) {

}

func (storage *Storage) Read(sessionId string) (Session, error) {

}

func (storage *Storage) Destroy(sessionId string) error {

}

func (storage *Storage) GarbageCollect(maxLifeTime int64) {

}

func (session *Session) Set(key, value interface{}) error {
	session.value[key] = value
	storage.Update(session.id)
}

func (session *Session) Get(key, interface{}) interface{} {
}

func (session *Session) Delete(key interface{}) error {
}

func (session *Session) Id() string {
}

func generateId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) Start(writer http.ReponseWriter, request *http.Request) (session Session) {
	manager.sessionMutex.Lock()
	defer manager.sessionMutex.Unlock()

	cookie, err := request.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		id := generateId()
		session, _ = manager.storage.Init(id)
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(id),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime),
		}

		http.SetCookie(writer, &cookie)
	} else {
		id, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.storage.Read(id)
	}

	return
}

func Wrapper(writer http.ResponseWriter, request *http.Request) {
	session := globalSessions.Start(writer, request)

}
