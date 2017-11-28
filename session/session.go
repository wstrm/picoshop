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

var storage *Storage

func init() {
	storage = &Storage{
		list: list.New(),
	}
}

type Manager struct {
	cookieName  string
	mutex       sync.Mutex
	maxLifeTime int64
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

func (store *Storage) Init(sessionId string) Session {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	value := make(map[interface{}]interface{}, 0)
	session := &Session{id: sessionId, lastAccess: time.Now(), value: value}
	element := store.list.PushBack(session)
	store.sessions[sessionId] = element

	return session
}

func (store *Storage) Read(sessionId string) (Session, error) {
	store.mutex.Lock() // do not defer Unlock() due to Storage.Init(sessionId)

	if element, ok := store.sessions[sessionsId]; ok {
		defer store.mutex.Unlock()
		return element.Value.(*Session), nil
	} else {
		store.mutex.Unlock() // store.Init locks mutex
		session := store.Init(sessionId)
		return session, nil
	}

	store.mutex.Unlock()
	return nil, nil
}

func (store *Storage) Destroy(sessionId string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if element, ok := store.sessions[sessionId]; ok {
		delete(store.sessions, sessionId)
		store.list.Remove(element)
		return nil
	}

	return nil
}

func (store *Storage) GarbageCollect(maxLifeTime int64) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for {
		element := store.list.Back()
		if element == nil {
			break
		}

		if (element.Value.(*Session).lastAccess.Unix() + maxLifeTime) < time.Now().Unix() {
			store.list.Remove(element)
			delete(store.sessions, element.Value.(*Session).id)
		} else {
			break
		}
	}
}

func (store *Storage) Update(sessionId string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if element, ok := store.sessions[sessionId]; ok {
		element.Value.(*Session).lastAccess = time.now()
		store.list.MoveToFront(element)
		return nil
	}

	return nil
}

func (session *Session) Set(key, value interface{}) {
	session.value[key] = value
	storage.Update(session.id)
}

func (session *Session) Get(key, interface{}) interface{} {
	storage.Update(session.id)

	if value, ok := session.value[key]; ok {
		return value
	} else {
		return nil
	}

	return nil
}

func (session *Session) Delete(key interface{}) error {
	delete(session.value, key)

	storage.Update(session.id)

	return nil
}

func (session *Session) Id() string {
	return session.id
}

func generateId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) Start(writer http.ReponseWriter, request *http.Request) (session Session) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	cookie, err := request.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		id := generateId()
		session, _ = storage.Init(id)
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
		session, _ = storage.Read(id)
	}

	return
}
