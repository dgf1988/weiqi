package weiqi

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	c_SessionCookieName = "WEIQI"
)

var (
	//锁
	sessionLocker sync.RWMutex
	//会话存储
	sessionStor map[string]*Session = make(map[string]*Session)
)

type Session struct {
	Id      string
	User    *User
	Timeout time.Time
}

func newSession(u *User) *Session {
	return &Session{
		sessionId(),
		u,
		time.Now().AddDate(0, 0, 1),
	}
}

func getSession(r *http.Request) *Session {
	s := &Session{}
	c, _ := r.Cookie(c_SessionCookieName)
	if c == nil {
		return nil
	}
	s.Id = c.Value

	sessionLocker.Lock()
	defer sessionLocker.Unlock()
	s_get, ok := sessionStor[s.Id]
	//id 相同， 并且user不为nil
	if !ok || s_get.Id != s.Id || s_get.User == nil {
		return nil
	}
	s.Timeout = s_get.Timeout
	s.User = s_get.User
	return s
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	sessionLocker.Lock()
	defer sessionLocker.Unlock()

	c, _ := r.Cookie(c_SessionCookieName)
	if c != nil {
		c.MaxAge = -1
		c.Expires = time.Now().AddDate(0, -1, 0)
		http.SetCookie(w, c)
	}

	_, ok := sessionStor[c.Value]
	if ok {
		delete(sessionStor, c.Value)
	}
}

func gcSession() int {
	sessionLocker.Lock()
	defer sessionLocker.Unlock()

	keys := make([]string, 0)
	for k, s := range sessionStor {
		if s.Timeout.Before(time.Now()) {
			keys = append(keys, k)
		}
	}
	for i := range keys {
		delete(sessionStor, keys[i])
	}
	return len(keys)
}

func sessionId() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	return getMd5(strconv.FormatInt(nano, 10) + strconv.FormatInt(rndNum, 10))
}

func (s *Session) Add(w http.ResponseWriter) {
	sessionLocker.Lock()
	defer sessionLocker.Unlock()

	http.SetCookie(w, &http.Cookie{Name: c_SessionCookieName, Value: s.Id, Expires: s.Timeout})
	sessionStor[s.Id] = s
}

//快速获取会话中的用户。
//并不是什么时候都需要持有一个会话对象。
func getSessionUser(r *http.Request) *User {
	s := getSession(r)
	if s != nil {
		return s.User
	}
	return nil
}
